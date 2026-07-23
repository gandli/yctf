# YCTF Changelog

> All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added

#### Documentation (Phase 0)
- **PRD.md**: Product requirements with 16 user stories (P0/P1/P2)
- **ARCHITECTURE.md**: System design, data flow diagrams, Docker network topology
- **DECISIONS.md**: 14 Architecture Decision Records (ADRs)
- **API.md**: 18 REST endpoints + WebSocket events + error format
- **DATABASE.md**: 8-table PostgreSQL schema + Redis data structures
- **DEVELOPMENT.md**: Dev setup, coding standards, testing strategy, Git workflow
- **SECURITY.md**: Threat model, flag/container/API security measures
- **TODO.md**: 6-phase roadmap (~12 weeks)
- **CONTRIBUTING.md**: Contribution guidelines
- **LICENSE**: AGPL-3.0 + trademark clause

#### Infrastructure (Phase 1)
- **docker-compose.yml**: PostgreSQL + Redis + Server + ClientApp
- **migrations/001_init.sql**: teams + users tables with indexes

#### Backend - TDD Cycles 1-8

| Cycle | Feature | Tests |
|-------|---------|-------|
| 1 | Health check endpoint | 1 |
| 2 | Flag generation (HMAC-SHA256) + validation | 4 |
| 3 | Password hashing (bcrypt cost=12) | 3 |
| 4 | JWT token generation & validation | 5 |
| 5 | User model with role/email/username validation | 5 |
| 6 | PostgreSQL connection pool (pgx) | 2 |
| 7 | Redis client + leaderboard + rate limit | 5 |
| 8 | Auth register handler | 7 |

**Total: 32 tests (22 pass + 5 DB-dependent skips)**

---

## [0.1.0] - 2026-07-23

### Added
- Project documentation suite
- Go module with chi, pgx, go-redis, jwt, bcrypt
- TDD foundation: 8 cycles, 32 tests
- Docker compose for local development
