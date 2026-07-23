# YCTF API Documentation

> RESTful API design and WebSocket event specification.

---

## Base Information

| Item | Value |
|------|-------|
| Base URL | `/api/v1` |
| Authentication | Bearer Token (JWT) |
| Format | JSON |
| CORS | Configurable (default `http://localhost:5173`) |

---

## Authentication

### POST `/auth/register`

Register a new user.

**Request:**
```json
{
  "username": "player1",
  "email": "player@example.com",
  "password": "SecurePass123!",
  "team_name": "TeamA",
  "invite_code": ""
}
```

**Response (201):**
```json
{
  "user": {
    "id": "uuid",
    "username": "player1",
    "email": "player@example.com",
    "role": "player",
    "team_id": "uuid"
  },
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

### POST `/auth/login`

**Request:**
```json
{
  "username": "player1",
  "password": "SecurePass123!"
}
```

**Response (200):**
```json
{
  "user": { "id": "uuid", "username": "player1", "role": "player" },
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

### POST `/auth/refresh`

**Request:**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response (200):**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

### POST `/auth/logout`

**Headers:** `Authorization: Bearer <token>`

**Response (204):** No Content

---

## Users

### GET `/users/me`

Get current user information.

**Response (200):**
```json
{
  "id": "uuid",
  "username": "player1",
  "email": "player@example.com",
  "role": "player",
  "team": {
    "id": "uuid",
    "name": "TeamA",
    "captain_id": "uuid",
    "score": 1500
  }
}
```

### PATCH `/users/me`

Update user info (nickname, password).

---

## Teams

### POST `/teams`

Create a team.

**Request:**
```json
{
  "name": "TeamA"
}
```

### POST `/teams/join`

Join a team via invite code.

**Request:**
```json
{
  "invite_code": "XXXXXX"
}
```

### GET `/teams/:id`

Get team details (members, solve history).

---

## Challenges

### GET `/challenges`

Get challenge list (filterable by category).

**Query Params:**
- `category` — web/pwn/crypto/re/misc/osint
- `page` — page number
- `limit` — per page

**Response (200):**
```json
{
  "challenges": [
    {
      "id": "uuid",
      "title": "SQL Injection 101",
      "category": "web",
      "points": 100,
      "solves": 15,
      "is_solved": false,
      "container_image": "yctf/web-sql1:latest"
    }
  ],
  "total": 42,
  "page": 1
}
```

### GET `/challenges/:id`

Get challenge details.

**Response (200):**
```json
{
  "id": "uuid",
  "title": "SQL Injection 101",
  "description": "...",
  "category": "web",
  "points": 100,
  "solves": 15,
  "is_solved": false,
  "container_image": "yctf/web-sql1:latest",
  "container_config": {
    "ports": [80],
    "env": ["FLAG_PREFIX=flag{"]
  },
  "attachments": [
    { "id": "uuid", "filename": "hint.txt", "url": "/api/v1/attachments/uuid" }
  ]
}
```

### POST `/challenges` (Admin/Author)

Create a challenge.

**Request:**
```json
{
  "title": "SQL Injection 101",
  "description": "...",
  "category": "web",
  "points": 100,
  "flag_template": "flag{...}",
  "container_image": "yctf/web-sql1:latest",
  "container_config": { "ports": [80] },
  "is_visible": true
}
```

---

## Container Instances

### POST `/challenges/:id/start`

Start a container instance.

**Response (201):**
```json
{
  "instance_id": "uuid",
  "container_id": "abc123",
  "host": "localhost",
  "port": 32001,
  "expires_at": "2026-07-23T12:00:00Z",
  "flag_env_injected": true
}
```

### POST `/instances/:id/stop`

Stop a container instance.

### GET `/instances/:id`

Get instance status.

**Response (200):**
```json
{
  "id": "uuid",
  "challenge_id": "uuid",
  "team_id": "uuid",
  "status": "running",
  "host": "localhost",
  "port": 32001,
  "started_at": "2026-07-23T10:00:00Z",
  "expires_at": "2026-07-23T12:00:00Z"
}
```

---

## Flag Submission

### POST `/submit`

Submit a flag.

**Request:**
```json
{
  "challenge_id": "uuid",
  "flag": "flag{sql_injection_success}"
}
```

**Response (200) — Correct:**
```json
{
  "correct": true,
  "message": "Flag correct! +100 points",
  "score_gained": 100,
  "team_score": 1600
}
```

**Response (200) — Incorrect:**
```json
{
  "correct": false,
  "message": "Incorrect flag"
}
```

**Response (429) — Rate limited:**
```json
{
  "error": "rate_limit_exceeded",
  "message": "Too many submissions. Try again in 30 seconds.",
  "retry_after": 30
}
```

---

## Scoreboard

### GET `/scoreboard`

Get the leaderboard.

**Response (200):**
```json
{
  "teams": [
    { "rank": 1, "team_id": "uuid", "team_name": "TeamA", "score": 2500, "solves": 8 },
    { "rank": 2, "team_id": "uuid", "team_name": "TeamB", "score": 2350, "solves": 7 }
  ],
  "updated_at": "2026-07-23T10:30:00Z"
}
```

### GET `/scoreboard/timeline`

Get team score timeline (for charts).

**Query Params:**
- `team_id` — optional, filter by team

**Response (200):**
```json
{
  "timeline": [
    { "timestamp": "2026-07-23T10:00:00Z", "score": 100, "event": "solved web-1" },
    { "timestamp": "2026-07-23T10:15:00Z", "score": 250, "event": "solved pwn-1" }
  ]
}
```

---

## Writeups

### POST `/writeups`

Submit a writeup.

**Request:**
```json
{
  "challenge_id": "uuid",
  "url": "https://example.com/writeup",
  "content": "optional markdown"
}
```

### GET `/writeups`

Get approved writeups list.

---

## Admin Endpoints

### GET `/admin/users`

User list (paginated, searchable).

### PATCH `/admin/users/:id`

Update user role/ban status.

### GET `/admin/submissions`

Submission audit log.

### GET `/admin/stats`

Platform statistics (users, submissions, containers).

---

## WebSocket Events

**Connection:** `ws://localhost:8080/ws`

**Handshake:** JWT token (query param or header)

### Event Types

| Event | Direction | Payload |
|-------|-----------|---------|
| `score_update` | Server → Client | `{ teams: [...], timestamp }` |
| `challenge_solved` | Server → Client | `{ team_name, challenge_title, points }` |
| `notification` | Server → Client | `{ type, message }` |
| `ping` | Bidirectional | `{}` |

### Example

```json
// Server → Client
{
  "event": "score_update",
  "data": {
    "teams": [
      { "rank": 1, "name": "TeamA", "score": 2500 }
    ],
    "timestamp": "2026-07-23T10:30:00Z"
  }
}
```

---

## Error Format

```json
{
  "error": "error_code",
  "message": "Human readable message",
  "details": {}
}
```

### Error Codes

| HTTP | Code | Description |
|------|------|-------------|
| 400 | `bad_request` | Invalid request format |
| 401 | `unauthorized` | Not authenticated |
| 403 | `forbidden` | Insufficient permissions |
| 404 | `not_found` | Resource not found |
| 409 | `conflict` | Resource conflict (e.g., duplicate registration) |
| 429 | `rate_limit_exceeded` | Too many requests |
| 500 | `internal_error` | Server error |
