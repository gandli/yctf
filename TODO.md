# YCTF Development TODO

> Milestones and task breakdown for the YCTF platform.

---

## Phase 1: Foundation (Week 1-2)

### Project Setup
- [ ] Initialize Go module (`go mod init`)
- [ ] Initialize React app (`npm create vite@latest clientapp -- --template react-ts`)
- [ ] Configure Mantine + Tailwind in clientapp
- [ ] Write `docker-compose.yml` (server + pg + redis + grafana)
- [ ] Configure CORS and chi router

### Database & Models
- [ ] Write initial migration (`001_init.sql`)
- [ ] Define Go structs for User, Team, Challenge, Instance, Submission
- [ ] Implement pgx connection pool
- [ ] Implement Redis connection

### Authentication
- [ ] Register endpoint (bcrypt password hash)
- [ ] Login endpoint (JWT access + refresh token)
- [ ] Middleware: JWT verification
- [ ] Middleware: RBAC (Admin, Author, Player)
- [ ] Frontend: Login page
- [ ] Frontend: Register page

---

## Phase 2: Challenge Management (Week 3-4)

### Challenge CRUD
- [ ] Create challenge (Admin/Author)
- [ ] List challenges (Player view, filtered by visibility)
- [ ] Update challenge
- [ ] Delete challenge (soft delete)
- [ ] Challenge categories (Web, PWN, Crypto, RE, Misc, OSINT)
- [ ] Frontend: Challenge list page
- [ ] Frontend: Challenge detail + submit form

### Flag Submission
- [ ] Submit flag endpoint
- [ ] Flag validation (HMAC comparison, constant-time)
- [ ] Rate limiting (Redis sliding window)
- [ ] Submission history per team
- [ ] Frontend: Submit form with feedback

---

## Phase 3: Container Distribution (Week 5-6)

### Docker Integration
- [ ] Docker SDK wrapper (`utils/docker.go`)
- [ ] Pull image if not cached
- [ ] Create container with env vars (unique flag)
- [ ] Start/stop container
- [ ] Container status tracking
- [ ] Port mapping (random host port → container port)
- [ ] Container GC (idle timeout)

### Instance Lifecycle
- [ ] Player clicks "Start Challenge" → creates instance
- [ ] Instance status page (running/stopped/expired)
- [ ] Auto-stop after competition end
- [ ] Frontend: Instance card with connection info
- [ ] Frontend: Start/stop buttons

---

## Phase 4: Scoreboard (Week 7-8)

### Real-time Leaderboard
- [ ] Redis Sorted Set for scores
- [ ] WebSocket hub (`ws/hub.go`)
- [ ] WS event on flag solved
- [ ] Score calculation (dynamic decay)
- [ ] Frontend: Live scoreboard table (Mantine DataTable)
- [ ] Frontend: Score timeline chart (Recharts or Mantine)

### Score Features
- [ ] First blood indicator
- [ ] Solve count per challenge
- [ ] Team profile with solve history
- [ ] Export scoreboard (CSV)

---

## Phase 5: Admin Panel (Week 9-10)

### Admin Dashboard
- [ ] Platform overview (users, teams, submissions, containers)
- [ ] User management (ban, promote, demote)
- [ ] Team management
- [ ] Challenge management (CRUD with container config)
- [ ] Submission audit log
- [ ] Grafana dashboard import (JSON)

### Writeup Module
- [ ] Writeup submission (URL or markdown)
- [ ] Writeup review queue
- [ ] Writeup scoring (Admin assigns bonus)

---

## Phase 6: Polish & Launch (Week 11-12)

### i18n
- [ ] Extract all UI strings
- [ ] Chinese translations
- [ ] Language switcher in navbar

### Security Hardening
- [ ] Container security (drop caps, read-only rootfs, no-new-privs)
- [ ] Content-Security-Policy headers
- [ ] SQL injection audit
- [ ] XSS audit
- [ ] Rate limiting audit

### Deployment & Docs
- [ ] GitHub Actions workflow (build + push to GHCR)
- [ ] Production `docker-compose.yml`
- [ ] Deployment guide
- [ ] User guide (Chinese + English)
- [ ] Write tests for critical paths (flag submit, auth, container lifecycle)

---

## Future Enhancements (Post-Launch)

- [ ] AWD mode support
- [ ] Flag rotation (periodic)
- [ ] Discord/Telegram webhook notifications
- [ ] Multi-site federation
- [ ] CTFtime integration
- [ ] CLI tool for challenge management
- [ ] Prometheus metrics

---

## Current Status

🚧 **Phase 1: Foundation** — Not started

Next action: Initialize Go module + React app, write initial docker-compose.yml
