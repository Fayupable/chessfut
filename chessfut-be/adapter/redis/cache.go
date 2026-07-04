package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	redislib "github.com/redis/go-redis/v9"

	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

const viewCounterWindow = 1 * time.Hour
const minCacheTTL = 10 * time.Minute

type Cache struct {
	client *redislib.Client
}

func (c *Cache) IncrementViewCount(ctx context.Context, username string) (int, error) {
	key := viewsKey(username)

	count, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		c.client.Expire(ctx, key, viewCounterWindow)
	}

	return int(count), nil
}
func (c *Cache) IsPromoted(ctx context.Context, username string) (bool, error) {
	exists, err := c.client.Exists(ctx, cardKey(username)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func NewCache(client *redislib.Client) *Cache {
	return &Cache{client: client}
}

var _ output.CachePort = (*Cache)(nil)

func cardKey(username string) string {
	return "card:" + username
}

func viewsKey(username string) string {
	return "views:" + username
}

func (c *Cache) GetCard(ctx context.Context, username string) (domain.Card, bool, error) {
	raw, err := c.client.Get(ctx, cardKey(username)).Bytes()
	if errors.Is(err, redislib.Nil) {
		return domain.Card{}, false, nil
	}
	if err != nil {
		return domain.Card{}, false, err
	}

	var model cardModel
	if err := json.Unmarshal(raw, &model); err != nil {
		return domain.Card{}, false, err
	}

	return fromCardModel(model), true, nil
}

func (c *Cache) SetCard(ctx context.Context, card domain.Card) error {
	data, err := json.Marshal(toCardModel(card))
	if err != nil {
		return err
	}

	ttl := time.Until(card.ExpiresAt)
	if ttl < minCacheTTL {
		ttl = minCacheTTL
	}

	return c.client.Set(ctx, cardKey(card.Player.Username), data, ttl).Err()
}
