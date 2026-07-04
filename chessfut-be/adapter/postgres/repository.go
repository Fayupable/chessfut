package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

type CardRepository struct {
	pool *pgxpool.Pool
}

func NewCardRepository(pool *pgxpool.Pool) *CardRepository {
	return &CardRepository{pool: pool}
}

var _ output.CardRepositoryPort = (*CardRepository)(nil)

func (r *CardRepository) Save(ctx context.Context, card domain.Card) error {
	data, err := json.Marshal(toCardDataModel(card))
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO player_cards (username, card_data, card_type, tier, computed_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (username) DO UPDATE SET
			card_data = EXCLUDED.card_data,
			card_type = EXCLUDED.card_type,
			tier = EXCLUDED.tier,
			computed_at = EXCLUDED.computed_at,
			expires_at = EXCLUDED.expires_at
	`

	_, err = r.pool.Exec(ctx, query,
		card.Player.Username, data, string(card.CardType), string(card.Tier),
		card.ComputedAt, card.ExpiresAt,
	)
	return err
}

func (r *CardRepository) FindByUsername(ctx context.Context, username string) (domain.Card, bool, error) {
	const query = `
		SELECT card_data, card_type, tier, computed_at, expires_at
		FROM player_cards
		WHERE username = $1
	`

	var (
		rawData             []byte
		cardType, tier      string
		computedAt, expires = domain.Card{}.ComputedAt, domain.Card{}.ExpiresAt
	)

	row := r.pool.QueryRow(ctx, query, username)
	if err := row.Scan(&rawData, &cardType, &tier, &computedAt, &expires); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Card{}, false, nil
		}
		return domain.Card{}, false, err
	}

	var model cardDataModel
	if err := json.Unmarshal(rawData, &model); err != nil {
		return domain.Card{}, false, err
	}

	return fromCardDataModel(model, cardType, tier, computedAt, expires), true, nil
}

func (r *CardRepository) FindStale(ctx context.Context, before time.Time, limit int) ([]domain.Card, error) {
	const query = `
		SELECT card_data, card_type, tier, computed_at, expires_at
		FROM player_cards
		WHERE expires_at < $1
		ORDER BY expires_at ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, before, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []domain.Card
	for rows.Next() {
		var (
			rawData        []byte
			cardType, tier string
			computedAt     time.Time
			expiresAt      time.Time
		)

		if err := rows.Scan(&rawData, &cardType, &tier, &computedAt, &expiresAt); err != nil {
			return nil, err
		}

		var model cardDataModel
		if err := json.Unmarshal(rawData, &model); err != nil {
			return nil, err
		}

		cards = append(cards, fromCardDataModel(model, cardType, tier, computedAt, expiresAt))
	}

	return cards, rows.Err()
}

func (r *CardRepository) FindTopByOVR(ctx context.Context, limit int) ([]domain.Card, error) {
	const query = `
		SELECT card_data, card_type, tier, computed_at, expires_at
		FROM player_cards
		ORDER BY (card_data->>'ovr')::int DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []domain.Card
	for rows.Next() {
		var (
			rawData        []byte
			cardType, tier string
			computedAt     time.Time
			expiresAt      time.Time
		)

		if err := rows.Scan(&rawData, &cardType, &tier, &computedAt, &expiresAt); err != nil {
			return nil, err
		}

		var model cardDataModel
		if err := json.Unmarshal(rawData, &model); err != nil {
			return nil, err
		}

		cards = append(cards, fromCardDataModel(model, cardType, tier, computedAt, expiresAt))
	}

	return cards, rows.Err()
}
