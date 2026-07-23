# YCTF Database Documentation

> PostgreSQL schema design and Redis data structures.

---

## ER Diagram (Core Entities)

```
+----------+      +--------------+      +-------------+
|  users   │      │   teams      │      │ challenges  |
+----------+      +--------------+      +-------------+
│ id (PK)  │--+   │ id (PK)      │      │ id (PK)     │
│ username │  └──>│ captain_id   │      │ title       │
│ email    │      │ name         │      │ category    │
│ password │      │ score        │      │ points      │
│ role     │      │ invite_code  │      │ is_visible  │
│ team_id  │      +--------------+      │ created_by  │
+----------+                            +-------------+
       │                                      │
       │       +--------------+              │
       │       │ submissions  │              │
       │       +--------------+              │
       │--┐    │ id (PK)      │    ┌--------+
          └──>│ user_id (FK)  │    │
               │ team_id (FK)  │<───┘
               │ challenge_id  │
               │ is_correct    │
               │ submitted_at  │
               └──────────────┘
```

---

## Table Definitions

### users

```sql
CREATE TABLE users (
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

CREATE INDEX idx_users_team ON users(team_id);
CREATE INDEX idx_users_email ON users(email);
```

### teams

```sql
CREATE TABLE teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(64) UNIQUE NOT NULL,
    captain_id      UUID REFERENCES users(id),
    score           INTEGER NOT NULL DEFAULT 0,
    invite_code     VARCHAR(16) UNIQUE,
    is_banned       BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### challenges

```sql
CREATE TABLE challenges (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title             VARCHAR(128) NOT NULL,
    description       TEXT NOT NULL DEFAULT '',
    category          VARCHAR(16) NOT NULL CHECK (category IN ('web', 'pwn', 'crypto', 're', 'misc', 'forensics', 'osint')),
    points            INTEGER NOT NULL DEFAULT 100,
    flag_template     VARCHAR(255),
    container_image   VARCHAR(255),
    container_config  JSONB DEFAULT '{}',
    is_visible        BOOLEAN NOT NULL DEFAULT FALSE,
    min_score_ratio   FLOAT DEFAULT 0.5,
    decay_threshold   INTEGER DEFAULT 100,
    created_by        UUID REFERENCES users(id),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_challenges_category ON challenges(category);
CREATE INDEX idx_challenges_visible ON challenges(is_visible);
```

### attachments

```sql
CREATE TABLE attachments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    challenge_id    UUID NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    filename        VARCHAR(255) NOT NULL,
    file_path       VARCHAR(500) NOT NULL,
    file_hash       VARCHAR(64),
    file_size       BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### instances

```sql
CREATE TABLE instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    challenge_id    UUID NOT NULL REFERENCES challenges(id),
    team_id         UUID NOT NULL REFERENCES teams(id),
    container_id    VARCHAR(64) UNIQUE,
    host_address    VARCHAR(255),
    host_port       INTEGER,
    status          VARCHAR(16) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'running', 'stopped', 'expired', 'error')),
    flag            VARCHAR(255) NOT NULL,
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ,
    last_heartbeat  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_instances_team ON instances(team_id);
CREATE INDEX idx_instances_status ON instances(status);
CREATE INDEX idx_instances_expires ON instances(expires_at);
```

### submissions

```sql
CREATE TABLE submissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    team_id         UUID NOT NULL REFERENCES teams(id),
    challenge_id    UUID NOT NULL REFERENCES challenges(id),
    flag_submitted  VARCHAR(255) NOT NULL,
    is_correct      BOOLEAN NOT NULL,
    ip_address      INET,
    submitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_submissions_team ON submissions(team_id);
CREATE INDEX idx_submissions_challenge ON submissions(challenge_id);
CREATE INDEX idx_submissions_user_time ON submissions(user_id, submitted_at);
```

### score_events

```sql
CREATE TABLE score_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID NOT NULL REFERENCES teams(id),
    challenge_id    UUID REFERENCES challenges(id),
    points_awarded  INTEGER NOT NULL,
    solved_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_score_events_team ON score_events(team_id);
CREATE INDEX idx_score_events_time ON score_events(solved_at);
```

### writeups

```sql
CREATE TABLE writeups (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    challenge_id    UUID NOT NULL REFERENCES challenges(id),
    team_id         UUID NOT NULL REFERENCES teams(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    url             VARCHAR(500),
    content         TEXT,
    is_approved     BOOLEAN DEFAULT NULL,
    score           INTEGER DEFAULT 0,
    reviewed_by     UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### games

```sql
CREATE TABLE games (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(128) NOT NULL,
    description     TEXT,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT FALSE,
    config          JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Database Migration

Use `golang-migrate` or `goose` for migration management.

```
migrations/
+-- 001_init.up.sql
+-- 001_init.down.sql
+-- 002_add_games.up.sql
+-- 002_add_games.down.sql
+-- ...
```

Migration commands:

```bash
# Migrate up
goose -dir migrations postgres "$DATABASE_URL" up

# Rollback one step
goose -dir migrations postgres "$DATABASE_URL" down

# Create new migration
goose -dir migrations create sql add_feature_x
```

---

## Redis Data Structures

### Leaderboard (Sorted Set)

```
Key: lb:scores
Type: ZSET
Member: team_id (string)
Score: integer points
```

### Rate Limiting (String + TTL)

```
Key: submit:limit:{user_id}
Type: STRING (counter)
TTL: 60 seconds
```

### Instance Cache (Hash)

```
Key: instance:{instance_id}
Type: HASH
Fields: challenge_id, team_id, container_id, host, port, status, flag, expires_at
TTL: 2 hours
```

### Challenge Solve Count (Set)

```
Key: challenge:{challenge_id}:solvers
Type: SET
Members: team_id
```

### Online Status (String)

```
Key: online:{team_id}
Type: STRING
Value: last_heartbeat timestamp
TTL: 5 minutes
```
