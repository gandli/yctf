-- migrations/001_init.up.sql

CREATE TABLE IF NOT EXISTS teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(64) UNIQUE NOT NULL,
    captain_id      UUID,
    score           INTEGER NOT NULL DEFAULT 0,
    invite_code     VARCHAR(16) UNIQUE,
    is_banned       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(32) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    role            VARCHAR(16) NOT NULL DEFAULT 'player' CHECK (role IN ('admin', 'author', 'player')),
    team_id         UUID REFERENCES teams(id) ON DELETE SET NULL,
    is_banned       BOOLEAN NOT NULL DEFAULT FALSE,
    last_login      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_team ON users(team_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
