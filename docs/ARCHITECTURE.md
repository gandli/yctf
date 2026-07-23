# YCTF Architecture

> System design, data flow, and module boundaries.

---

## 1. System Overview

```
+-----------------------------------------------------------------------+
|                         Client (Browser)                               |
|  React SPA + WebSocket Client + i18n                                  |
+--------------------------------+--------------------------------------+
                                 | HTTP / WSS
+--------------------------------v--------------------------------------+
|                          Nginx / CORS                                  |
+----------+-----------------------------------------------------------+
           |
+----------v-----------------------------------------------------------+
|                        Go Backend (chi)                                |
|  +----------+----------+----------+----------+----------+             |
|  |  Auth    | Challenge|  Submit  | Container|  Admin   |             |
|  |Controller|Controller|Controller|Controller|Controller|             |
|  +----------+----------+----------+----------+----------+             |
|  +-----------------------------------------------------+             |
|  |              Middleware (RBAC, Rate Limit)          |             |
|  +-----------------------------------------------------+             |
+------+--------------+--------------+----------------------------------+
       |              |              |
+------v------+ +-----v------+ +----v------+
| PostgreSQL  │ |   Redis    │ |  Docker   │
│ (Primary)   │ | (Cache/RT) │ | (Daemon)  |
+-------------+ +------------+ +-----------+
```

---

## 2. Module Boundaries

### 2.1 Backend (Go)

```
src/
+-- cmd/server/          # Entry point, graceful shutdown
+-- controllers/         # HTTP handlers, request validation
|   +-- auth.go          # Login/register/logout
|   +-- challenge.go     # CRUD challenges, listing
|   +-- submit.go        # Flag submission, validation
|   +-- container.go     # Instance start/stop/status
|   +-- scoreboard.go    # Leaderboard query, timeline
|   +-- admin.go         # User/challenge management
|   +-- writeup.go       # Writeup submission
+-- db/
|   +-- models/          # GORM-free struct definitions
|   +-- migrations/      # SQL migration files
|   +-- queries/         # Raw SQL queries
+-- middleware/
|   +-- auth.go          # JWT verification
|   +-- rbac.go          # Role-based access control
|   +-- ratelimit.go     # Redis-backed rate limiting
|   +-- cors.go          # CORS configuration
+-- tasks/
|   +-- flag_rotator.go  # Periodic flag rotation
|   +-- container_gc.go  # Idle container cleanup
|   +-- score_sync.go    # PG -> Redis scoreboard sync
+-- utils/
|   +-- flag.go          # Flag generation & validation
|   +-- hash.go          # Password hashing (bcrypt)
|   +-- docker.go        # Docker SDK wrapper
+-- ws/
    +-- hub.go           # WebSocket connection hub
    +-- client.go        # WS client abstraction
    +-- message.go       # Event types (score, notify)
```

### 2.2 Frontend (React)

```
clientapp/src/
+-- app/
|   +-- routes.tsx       # Route definitions
|   +-- stores/          # Zustand stores (auth, challenge, score)
|   +-- pages/
|       +-- Home.tsx     # Landing + CTF info
|       +-- Board.tsx    # Real-time leaderboard
|       +-- Challenges.tsx  # Challenge list + detail
|       +-- Profile.tsx  # Team profile, settings
|       +-- Admin/       # Admin panel
|           +-- Dashboard.tsx
|           +-- Challenges.tsx
|           +-- Users.tsx
|           +-- Monitor.tsx
+-- components/
|   +-- ui/              # Mantine-based primitives
|   +-- layout/          # Navbar, footer, sidebar
|   +-- challenge/       # Challenge card, submit form, instance panel
|   +-- scoreboard/      # Scoreboard table, timeline chart
+-- hooks/
|   +-- useWebSocket.ts  # WS connection & auto-reconnect
|   +-- useAuth.ts       # Auth state management
|   +-- useChallenge.ts  # Challenge data fetching
+-- utils/
    +-- api.ts           # Axios instance with interceptors
    +-- i18n.ts          # react-i18next config
    +-- flag.ts          # Flag format helpers
```

---

## 3. Data Flow

### 3.1 Flag Submission

```
Player submits flag
        |
        v
+-- Frontend ----------------------------------+
| 1. Trim, validate format                    |
| 2. POST /api/submit { challenge_id, flag } |
+------------------------+--------------------+
                         |
+-- Backend -------------v--------------------+
| 3. Rate limit check (Redis)                 |
| 4. JWT -> identify user & team              |
| 5. Lookup challenge flag in Redis/PG        |
| 6. Compare (constant-time)                  |
| 7. If correct:                              |
|    a. Write submission record to PG         |
|    b. Update team score in Redis ZSET       |
|    c. Broadcast WS event (score update)     |
| 8. Return result                            |
+---------------------------------------------+
```

### 3.2 Container Lifecycle

```
Admin creates challenge with Docker image
        |
        v
+-- Admin -------------------------------------+
| 1. Define image, ports, env vars            |
| 2. Set flag template                        |
| 3. Challenge saved to PG                    |
+------------------------+--------------------+
                         | Player clicks "Start"
+-- Backend -------------v--------------------+
| 4. Pull image if not cached                 |
| 5. Generate unique flag (team+challenge)    |
| 6. Create container with env injection      |
| 7. Store instance mapping in PG+Redis       |
| 8. Return connection info to player         |
| 9. Background GC after idle timeout         |
+---------------------------------------------+
```

### 3.3 Real-time Scoreboard

```
Flag solved event
        |
        v
+-- Backend ----------------------------------+
| 1. Update team score in Redis ZSET          |
| 2. Publish event to WS hub                 |
+------------------------+--------------------+
                         |
+-- WebSocket Hub -------v-------------------+
| 3. Fan-out to all connected clients        |
+------------------------+--------------------+
                         |
+-- Frontend ------------v-------------------+
| 4. Mantine Table re-render                 |
| 5. Toast notification (optional)           |
+--------------------------------------------+
```

---

## 4. Database Schema (Core)

```sql
-- Users & Teams
users (id, username, email, password_hash, role, team_id, created_at)
teams (id, name, captain_id, score, invited_code, created_at)

-- Challenges
challenges (id, title, description, category, points, flag_template,
            container_image, container_config, is_visible, created_by, created_at)
attachments (id, challenge_id, filename, url, hash)

-- Container Instances
instances (id, challenge_id, team_id, container_id, internal_port,
           status, started_at, expires_at, flag)

-- Submissions
submissions (id, user_id, team_id, challenge_id, flag_submitted,
             is_correct, ip_address, submitted_at)

-- Score timeline (for charts)
score_events (id, team_id, challenge_id, points_awarded, solved_at)

-- Writeups
writeups (id, challenge_id, team_id, user_id, url, content,
          is_approved, score, created_at)
```

---

## 5. Redis Data Structures

| Key Pattern | Type | Purpose |
|-------------|------|---------|
| `lb:scores` | Sorted Set | Global leaderboard (member=team_id, score=points) |
| `submit:limit:{user_id}` | String (counter) | Rate limit flag submissions |
| `instance:{id}` | Hash | Container instance metadata |
| `challenge:{id}:sponsors` | Set | Teams who solved (for dynamic scoring calc) |
| `online:{team_id}` | String | Last heartbeat (presence) |

---

## 6. Docker Network Architecture

```
+---------------------------------------------+
|           yctf_network                      |
|                                             |
|  +----------+  +----------+  +-------+     |
|  |  yctf-   |  |  yctf-   |  | yctf- |     |
|  |  server  |  |  redis   |  |  pg   |     |
|  +----------+  +----------+  +-------+     |
|                                             |
|  +-------------------------------------+   |
|  |    challenge containers             |   |
|  |  +-----+ +-----+ +-----+            |   |
|  |  |web-1| |pwn-1| |cry-1|  ...      |   |
|  |  +-----+ +-----+ +-----+            |   |
|  +-------------------------------------+   |
+---------------------------------------------+
```

- Challenge containers join the same Docker network
- Backend communicates with containers via internal Docker DNS
- Player accesses containers via published ports or internal IP

---

## 7. Security Considerations

| Threat | Mitigation |
|--------|-----------|
| Flag sharing | Per-team unique flags via container env injection |
| Brute force flag | Redis rate limit (10 submissions/min per user) |
| Container escape | Drop capabilities, read-only rootfs, no-new-privileges |
| SQL injection | Parameterized queries (pgx native) |
| XSS | React auto-escape + Content-Security-Policy headers |
| WS hijack | JWT token in WS handshake |

---

## 8. Scaling Path

| Stage | Action |
|-------|--------|
| < 50 teams | Docker compose single node (current) |
| 50-200 teams | Add Redis Sentinel, separate PG server |
| 200+ teams | Docker Swarm / K3s, load balancer |
| Multi-site | Federation via API, shared Redis for cross-site LB |
