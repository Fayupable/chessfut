CREATE TABLE player_cards (
                              username     VARCHAR(50) PRIMARY KEY,
                              card_data    JSONB       NOT NULL,
                              card_type    VARCHAR(10) NOT NULL,
                              tier         VARCHAR(10) NOT NULL,
                              computed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                              expires_at   TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_player_cards_expires_at ON player_cards (expires_at);
CREATE INDEX idx_player_cards_tier ON player_cards (tier);