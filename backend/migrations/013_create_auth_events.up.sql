CREATE TABLE auth_events (
    id          UUID PRIMARY KEY,
    event_type  VARCHAR(50)  NOT NULL,
    user_id     UUID         REFERENCES users(id) ON DELETE SET NULL,
    email       VARCHAR(255),
    ip_address  VARCHAR(45),
    user_agent  TEXT,
    metadata    JSONB        NOT NULL DEFAULT '{}',
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_events_created_at  ON auth_events(created_at DESC);
CREATE INDEX idx_auth_events_email       ON auth_events(email);
CREATE INDEX idx_auth_events_user_id     ON auth_events(user_id);
CREATE INDEX idx_auth_events_event_type  ON auth_events(event_type);
