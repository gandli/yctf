# Technical Decisions

> Architecture Decision Records (ADRs) for YCTF.

---

## ADR-001: Competition Mode — Jeopardy

**Status:** Accepted

**Context:** The platform needs to support a competition format. Options are Jeopardy (fixed challenges), AWD (attack-defense with per-team targets), or hybrid.

**Decision:** Start with Jeopardy only. AWD adds significant complexity (per-team VMs, network isolation, traffic capture, automated scoring of patching). Jeopardy covers the core CTF experience and is the most common format for training and small-to-medium competitions.

**Consequences:**
- Simple container model (shared image, per-team flag injection)
- No need for complex network topologies
- AWD can be added later as a separate module if needed

---

## ADR-002: Challenge Distribution — Dynamic Containers

**Status:** Accepted

**Context:** Challenges can be static (shared attachment/URL) or dynamic (per-team container with unique flag).

**Decision:** Fully dynamic container distribution. Each team gets an isolated container instance with a unique flag injected via environment variables.

**Consequences:**
- Eliminates flag sharing between teams
- Requires container orchestration (Docker)
- Higher resource usage per team
- Need container lifecycle management (start, stop, GC)

---

## ADR-003: Container Orchestration — Docker Single Node

**Status:** Accepted

**Context:** Need to run challenge containers. Options: Docker single, Swarm, Kubernetes, K3s.

**Decision:** Docker single node with docker-compose. No orchestration layer.

**Consequences:**
- Simplest deployment (single binary + compose file)
- Can handle 50-100 teams on modest hardware
- No auto-healing or load balancing
- Migration path: abstract container operations behind an `Instancer` interface for future Swarm/K8s support

---

## ADR-004: Backend Language — Go

**Status:** Accepted

**Context:** CTF platform backend needs high concurrency for flag submissions, container management, and WebSocket connections.

**Decision:** Go 1.24 with chi router.

**Consequences:**
- Single binary deployment
- Excellent concurrency (goroutines for WebSocket fan-out)
- Strong Docker SDK support
- Good PostgreSQL drivers (pgx)
- Fast compilation and execution

---

## ADR-005: Frontend — React + Vite + Mantine + Tailwind

**Status:** Accepted

**Context:** Need a modern, responsive UI for scoreboard, challenge listing, admin panel.

**Decision:** React 18 with Vite build tool, Mantine component library, Tailwind CSS.

**Consequences:**
- Mantine provides rich admin components (Table, Tabs, Notifications)
- Tailwind for utility-first styling
- Vite for fast HMR in development
- Real-time updates via socket.io-client
- TypeScript for type safety

---

## ADR-006: Database — PostgreSQL

**Status:** Accepted

**Context:** Need persistent storage for users, teams, challenges, submissions, and score history.

**Decision:** PostgreSQL as primary database.

**Consequences:**
- ACID compliance for score transactions
- JSONB for flexible config storage (challenge flags, container config)
- Window functions for leaderboard queries
- Proven in CTF platforms (GZCTF, rCTF both use PG)

---

## ADR-007: Cache & Realtime — Redis

**Status:** Accepted

**Context:** Leaderboard needs real-time updates and flag submission needs rate limiting.

**Decision:** Redis for caching, rate limiting, and real-time leaderboard (Sorted Sets).

**Consequences:**
- Sorted Sets provide O(log N) leaderboard operations
- Rate limiting via sliding window counters
- Session storage if needed
- Pub/Sub for WebSocket fan-out across multiple backend instances

---

## ADR-008: Flag Strategy — Container Environment Injection

**Status:** Accepted

**Context:** Need to generate unique flags per team per challenge.

**Decision:** Generate flag as `flag{HMAC(team_id + challenge_id + secret)}` and inject via container environment variable at startup.

**Consequences:**
- Flags are unique per team (no sharing)
- No need to bake flags into images
- Flag rotation requires container restart
- HMAC prevents flag forgery

---

## ADR-009: Authentication — JWT with RBAC

**Status:** Accepted

**Context:** Need user authentication and role-based access control.

**Decision:** JWT tokens with three roles (Admin, Author, Player).

**Consequences:**
- Stateless authentication (easy horizontal scaling)
- Roles enforced via middleware
- Token refresh via secure httpOnly cookies
- Initial implementation: access token (15min) + refresh token (7d)

---

## ADR-010: Scoring — Dynamic Decay

**Status:** Accepted

**Context:** Static scoring (fixed points per challenge) is simple but doesn't reward early solves.

**Decision:** Dynamic scoring with decay formula: `score = base_points * (1 - solvers/decay_threshold)`. Minimum floor at 50% of base.

**Consequences:**
- Rewards solving before others
- Prevents farming easy points late in competition
- Formula configurable per challenge
- Later upgrade path: exponential decay curve (like GZCTF)

---

## ADR-011: Internationalization — i18n

**Status:** Accepted

**Context:** Platform needs to support Chinese and English.

**Decision:** React-i18next for frontend, Go `gettext` or similar for backend errors.

**Consequences:**
- All UI strings externalized
- Language auto-detection from browser, override in user settings
- Translation files in `/i18n/{lang}/`

---

## ADR-012: Monitoring — Grafana

**Status:** Accepted

**Context:** Need visibility into platform health, container usage, submission patterns.

**Decision:** Grafana dashboards with PostgreSQL as data source.

**Consequences:**
- Pre-built dashboards for common CTF metrics
- No additional data pipeline needed (direct PG queries)
- Future: add Prometheus for container metrics

---

## ADR-013: Deployment — Docker Compose + GitHub Actions

**Status:** Accepted

**Context:** Need simple deployment and CI/CD.

**Decision:** Docker Compose for runtime, GitHub Actions for build/push to GHCR.

**Consequences:**
- Single `docker compose up -d` for production
- Automated image builds on push to main
- No complex Kubernetes manifests
- Offline deployment possible via image export

---

## ADR-014: Project Structure — A1CTF-Inspired

**Status:** Accepted

**Context:** Need clear separation of concerns between frontend and backend.

**Decision:** Follow A1CTF structure: `clientapp/` for React, `src/` for Go, `migrations/`, `i18n/`, `grafana/`.

**Consequences:**
- Familiar to contributors from similar projects
- Clear module boundaries
- Independent versioning of frontend and backend
