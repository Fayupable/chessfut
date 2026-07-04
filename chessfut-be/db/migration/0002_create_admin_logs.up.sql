CREATE TABLE admin_logs (
                            id               SERIAL PRIMARY KEY,
                            action           VARCHAR(50) NOT NULL,
                            target_username  VARCHAR(50),
                            created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                            details          JSONB
);

CREATE INDEX idx_admin_logs_action ON admin_logs (action);
CREATE INDEX idx_admin_logs_created_at ON admin_logs (created_at);