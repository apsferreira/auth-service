-- Migration 014: OAuth identities for social login (Google, future providers)

CREATE TABLE IF NOT EXISTS oauth_identities (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider    VARCHAR(50)  NOT NULL,       -- 'google', 'github', etc.
    provider_id VARCHAR(255) NOT NULL,       -- Google's 'sub' claim
    email       VARCHAR(255) NOT NULL,
    avatar_url  TEXT,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_id)
);

CREATE INDEX IF NOT EXISTS idx_oauth_identities_user_id    ON oauth_identities(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_identities_provider   ON oauth_identities(provider, provider_id);
