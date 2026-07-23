# Product Requirements Document (PRD)

> YCTF Platform v1.0

---

## 1. Product Overview

**YCTF** is an open-source Jeopardy-style CTF (Capture The Flag) competition platform designed for security teams, CTF organizers, and training environments.

### Mission

Enable anyone to host professional-grade CTF competitions with minimal infrastructure — one command to deploy, intuitive UI for management, and rock-solid flag security.

### Target Users

| User | Need |
|------|------|
| **CTF Organizer** | Create competitions, manage challenges, monitor progress |
| **Challenge Author** | Design and publish challenges with dynamic containers |
| **Player** | Solve challenges, track progress, compete on leaderboard |
| **Admin** | Platform configuration, user management, monitoring |

---

## 2. Product Goals

| Goal | Metric |
|------|--------|
| Simplicity | Deploy with `docker compose up -d` in < 5 minutes |
| Security | Zero flag-sharing incidents, rate-limited submissions |
| Real-time | Leaderboard updates < 1s after submission |
| Scalability | Support 50-100 teams on single node |
| International | Chinese + English from day one |

---

## 3. Scope (v1.0)

### In Scope

- Jeopardy-style competition format
- Dynamic Docker container challenges
- Per-team unique flag injection
- Real-time WebSocket leaderboard
- Dynamic scoring with decay
- Three-tier RBAC (Admin, Author, Player)
- Challenge categories (Web, PWN, Crypto, RE, Misc, OSINT)
- Team creation and management
- Writeup submission (URL-based)
- Grafana monitoring dashboards
- i18n (Chinese + English)

### Out of Scope (v1.0)

- AWD (Attack with Defense) mode
- Automated challenge evaluation
- Mobile native app
- Payment/registration fee integration
- CTFtime official integration
- Flag rotation during competition

---

## 4. User Stories

### Organizer

| ID | Story | Priority |
|----|-------|----------|
| ORG-1 | As an organizer, I want to create a competition with start/end times so that all players compete within the same window | P0 |
| ORG-2 | As an organizer, I want to see a real-time dashboard of platform health so that I can detect issues early | P0 |
| ORG-3 | As an organizer, I want to export final scores as CSV so that I can publish results | P1 |

### Challenge Author

| ID | Story | Priority |
|----|-------|----------|
| AUT-1 | As an author, I want to create a challenge with a Docker image so that each player gets an isolated instance | P0 |
| AUT-2 | As an author, I want to configure the flag template so that unique flags are generated per team | P0 |
| AUT-3 | As an author, I want to attach static files (for non-container challenges) so that players can download them | P1 |

### Player

| ID | Story | Priority |
|----|-------|----------|
| PLA-1 | As a player, I want to see all active challenges categorized by type so that I can choose what to solve | P0 |
| PLA-2 | As a player, I want to submit a flag and get instant feedback so that I know if I solved it | P0 |
| PLA-3 | As a player, I want to start a container instance with one click so that I can access the challenge environment | P0 |
| PLA-4 | As a player, I want to view the live leaderboard so that I can track my team's rank | P0 |
| PLA-5 | As a player, I want to join or create a team so that I can compete with others | P0 |
| PLA-6 | As a player, I want to submit a writeup after solving a challenge so that I can share my approach | P2 |

### Admin

| ID | Story | Priority |
|----|-------|----------|
| ADM-1 | As an admin, I want to ban a user so that I can remove bad actors | P0 |
| ADM-2 | As an admin, I want to review writeups so that I can award bonus points | P1 |
| ADM-3 | As an admin, I want to monitor container usage so that I can manage resources | P1 |

---

## 5. Functional Requirements

### 5.1 Authentication & Authorization

- User registration with email verification
- Login with JWT (access token + refresh token)
- Three roles: Admin, Author, Player
- Optional team review mode (admin approves registrations)
- Password reset via email (SMTP)

### 5.2 Challenge Management

- Create/edit/delete challenges
- Categories: Web, PWN, Crypto, RE, Misc, Forensics, OSINT
- Challenge types:
  - **Dynamic Container**: Docker image + env vars for flag injection
  - **Static Attachment**: File download
- Visibility control (visible/hidden)
- Point value configuration
- Dynamic scoring parameters (decay threshold)

### 5.3 Container Management

- Docker image pull and cache
- Container creation with unique flag injection
- Port mapping (auto-assign host port)
- Container status monitoring (running/stopped/expired)
- Auto-GC after idle timeout or competition end

### 5.4 Flag Submission

- Submit flag for a challenge
- Instant validation (correct/incorrect)
- Rate limiting: max 10 submissions/minute per user
- Submission history per team
- Anti-sharing: unique flags per team

### 5.5 Scoreboard

- Real-time ranking via WebSocket
- Team score timeline (chart)
- First-blood bonus display
- Challenge solve count
- Export to CSV

### 5.6 Team Management

- Create team (captain)
- Join team (invite code)
- Team profile with solve history
- Team score aggregation

### 5.7 Writeup System

- Submit writeup URL after solving
- Admin review queue
- Approval/rejection workflow

---

## 6. Non-Functional Requirements

| Category | Requirement |
|----------|-------------|
| Performance | Flag submission response < 200ms (p95) |
| Availability | 99.9% uptime during competition |
| Security | OWASP Top 10 mitigation, container isolation |
| Usability | Mobile-responsive UI, accessible (WCAG 2.1 AA) |
| Maintainability | Test coverage > 70% for core paths |
| Portability | Runs on any Docker-compatible Linux host |

---

## 7. Technical Constraints

- Go 1.24+ (backend)
- React 18+ (frontend)
- PostgreSQL 15+ (database)
- Redis 7+ (cache)
- Docker 24+ (runtime)

---

## 8. Release Criteria

| Check | Status |
|-------|--------|
| All P0 user stories implemented | ⬜ |
| All P0 user stories tested (TDD) | ⬜ |
| Documentation complete | ✅ |
| Security audit passed | ⬜ |
| Performance benchmark met | ⬜ |
| Docker compose deployment verified | ⬜ |
| i18n verified (ZH + EN) | ⬜ |

---

## 9. Success Metrics

| Metric | Target |
|--------|--------|
| Time to deploy (new instance) | < 5 minutes |
| Concurrent teams supported | 50-100 |
| Flag submission latency (p95) | < 200ms |
| Leaderboard update latency | < 1s |
| Test coverage (core paths) | > 70% |
