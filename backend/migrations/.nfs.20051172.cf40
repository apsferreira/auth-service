-- Migration 011: Add username and password_hash columns for admin panel authentication
-- Admin users authenticate with username/password; consumer apps continue using OTP

ALTER TABLE users ADD COLUMN username VARCHAR(100);
ALTER TABLE users ADD COLUMN password_hash VARCHAR(255);

-- Partial unique index: only enforces uniqueness when username is not NULL
CREATE UNIQUE INDEX idx_users_username ON users(username) WHERE username IS NOT NULL;
