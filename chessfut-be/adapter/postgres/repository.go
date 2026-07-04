package postgres

import (
	"context"
	"encoding/json"
	"errors"

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
