package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fayupable/chessfut-be/adapter/health"
	"github.com/jackc/pgx/v5/pgxpool"
	redislib "github.com/redis/go-redis/v9"

	"github.com/fayupable/chessfut-be/adapter/chesscom"
	httpadapter "github.com/fayupable/chessfut-be/adapter/http"
	"github.com/fayupable/chessfut-be/adapter/postgres"
	redisadapter "github.com/fayupable/chessfut-be/adapter/redis"
	"github.com/fayupable/chessfut-be/application/service"
	"github.com/fayupable/chessfut-be/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()
	ctx := context.Background()

	if err := postgres.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied successfully")

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	redisClient := redislib.NewClient(&redislib.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	defer redisClient.Close()

	chessComClient := chesscom.NewClient(cfg.ChessComUserAgent, cfg.ChessComBaseURL)
	cardRepository := postgres.NewCardRepository(pool)
	cache := redisadapter.NewCache(redisClient)

	postgresChecker := health.NewChecker("postgres", func(ctx context.Context) error {
		return pool.Ping(ctx)
	})
	redisChecker := health.NewChecker("redis", func(ctx context.Context) error {
		return redisClient.Ping(ctx).Err()
	})
	postgresChecker.Start(ctx)
	redisChecker.Start(ctx)

	getFastCard := service.NewGetFastCardService(chessComClient, cardRepository, cache)
	getDetailedCard := service.NewGetDetailedCardService(chessComClient, cardRepository, cache)
	getLeaderboard := service.NewGetLeaderboardService(cardRepository)
	getStats := service.NewGetStatsService(cardRepository, cache)
	searchPlayers := service.NewSearchPlayersService(cardRepository)
	refreshCard := service.NewRefreshCardService(chessComClient, cardRepository, cache)
	syncTitledPlayers := service.NewSyncTitledPlayersService(chessComClient, cardRepository)
	refreshStaleCards := service.NewRefreshStaleCardsService(chessComClient, cardRepository, cache)

	cardHandler := httpadapter.NewCardHandler(getFastCard, getDetailedCard, getLeaderboard, getStats, searchPlayers)
	adminHandler := httpadapter.NewAdminHandler(refreshCard, syncTitledPlayers, refreshStaleCards)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/player/{username}", httpadapter.CORSMiddleware(httpadapter.LoggingMiddleware(cardHandler.GetFastCard)))
	mux.HandleFunc("GET /api/v1/player/{username}/detailed", httpadapter.CORSMiddleware(httpadapter.LoggingMiddleware(cardHandler.GetDetailedCard)))
	mux.HandleFunc("GET /api/v1/leaderboard", httpadapter.CORSMiddleware(httpadapter.LoggingMiddleware(cardHandler.GetLeaderboard)))
	mux.HandleFunc("GET /api/v1/stats", httpadapter.CORSMiddleware(httpadapter.LoggingMiddleware(cardHandler.GetStats)))
	mux.HandleFunc("GET /api/v1/search", httpadapter.CORSMiddleware(httpadapter.LoggingMiddleware(cardHandler.SearchPlayers)))

	mux.HandleFunc("POST /api/admin/refresh/{username}", httpadapter.LoggingMiddleware(httpadapter.AdminAuthMiddleware(cfg.AdminAPIKey, adminHandler.RefreshCard)))
	mux.HandleFunc("POST /api/admin/sync-titled", httpadapter.LoggingMiddleware(httpadapter.AdminAuthMiddleware(cfg.AdminAPIKey, adminHandler.SyncTitledPlayers)))
	mux.HandleFunc("POST /api/admin/refresh-stale", httpadapter.LoggingMiddleware(httpadapter.AdminAuthMiddleware(cfg.AdminAPIKey, adminHandler.RefreshStaleCards)))

	startBackgroundJobs(refreshStaleCards, syncTitledPlayers)

	log.Printf("chessfut-be listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func startBackgroundJobs(refreshStaleCards *service.RefreshStaleCardsService, syncTitledPlayers *service.SyncTitledPlayersService) {
	staleTicker := time.NewTicker(10 * time.Minute)
	go func() {
		for range staleTicker.C {
			count, err := refreshStaleCards.Execute(context.Background(), 50)
			if err != nil {
				log.Printf("refresh stale cards failed: %v", err)
				continue
			}
			log.Printf("refreshed %d stale cards", count)
		}
	}()

	syncTicker := time.NewTicker(24 * time.Hour)
	go func() {
		for range syncTicker.C {
			count, err := syncTitledPlayers.Execute(context.Background())
			if err != nil {
				log.Printf("sync titled players failed: %v", err)
				continue
			}
			log.Printf("synced %d new titled players", count)
		}
	}()
}
