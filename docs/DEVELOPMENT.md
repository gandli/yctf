# YCTF Developer Documentation

> Development environment setup, coding standards, testing strategy, and deployment guide.

---

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Coding Standards](#coding-standards)
- [Testing Strategy](#testing-strategy)
- [Git Workflow](#git-workflow)
- [Release Process](#release-process)

---

## Development Environment Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.24+ | Backend development |
| Node.js | 20+ | Frontend development |
| Docker | 24+ | Container runtime |
| Docker Compose | v2+ | Local development environment |
| Make | any | Task orchestration |

### Quick Start

```bash
# Clone
git clone https://github.com/gandli/yctf.git && cd yctf

# Start dependencies (PostgreSQL + Redis)
docker compose -f compose.dev.yml up -d

# Backend (new terminal)
cd src
go mod tidy
go run cmd/server/main.go

# Frontend (new terminal)
cd clientapp
npm install
npm run dev
```

### Access

| Service | URL |
|---------|-----|
| Frontend Dev Server | http://localhost:5173 |
| Backend API | http://localhost:8080 |
| Grafana | http://localhost:3001 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

### Environment Variables

```bash
# Backend (.env)
DATABASE_URL=postgres://yctf:yctf@localhost:5432/yctf?sslmode=disable
REDIS_URL=redis://localhost:6379/0
JWT_SECRET=dev-secret-change-in-production
PORT=8080
CORS_ORIGINS=http://localhost:5173

# Frontend (.env)
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

---

## Coding Standards

### Go Standards

| Rule | Description |
|------|-------------|
| Formatting | `gofmt -w` (enforced by CI) |
| Linting | `golangci-lint run` (`.golangci.yml`) |
| Naming | PascalCase for exported, camelCase for unexported, all-caps for acronyms (HTTP, ID) |
| Error Handling | `if err != nil` must be handled, never `_ = err` |
| Interfaces | Single responsibility, file named `interface.go` |
| Tests | Same directory `*_test.go`, table-driven tests preferred |
| Dependency Injection | Via interfaces, no global variables (except `const`) |

### React/TypeScript Standards

| Rule | Description |
|------|-------------|
| Components | Function components + Hooks, no class components |
| File Naming | PascalCase for components, camelCase for utilities |
| Types | Strict mode, no `any` (except untyped third-party libs) |
| Props Interface | `interface ComponentNameProps` |
| Styling | Tailwind first, Mantine `sx` second |
| Imports | Absolute path `@/` → `src/` |
| State | Zustand stores, avoid prop drilling |

### Commit Convention

```
feat:     New feature
fix:      Bug fix
docs:     Documentation change
style:    Code formatting (no functional change)
refactor: Code restructuring
perf:     Performance improvement
test:     Test-related changes
chore:    Build/dependency/CI configuration
ci:       CI configuration
```

Example: `feat: add flag submission rate limiting`

---

## Testing Strategy

### TDD Iron Law

> **No production code without a failing test first.**

```bash
# Backend tests
cd src
go test ./... -v -count=1              # Full suite
go test ./controllers/ -run TestSubmit  # Specific test

# Frontend tests
cd clientapp
npm run test                            # Vitest unit tests
npm run test:e2e                        # Playwright E2E
```

### Testing Pyramid

```
        /  E2E (Playwright)  \         ← User journeys
       / Integration (API Test) \       ← Module collaboration
      /    Unit (Go test / Vitest) \    ← Foundation units
```

### Coverage Targets

| Module | Target |
|--------|--------|
| `utils/` | 90%+ |
| `controllers/` | 80%+ |
| `db/queries/` | 85%+ |
| `middleware/` | 75%+ |
| Frontend hooks | 80%+ |

### Critical Test Scenarios

1. **Flag Submission**: correct flag → score; incorrect flag → rejected; rate limit → 429
2. **Container Lifecycle**: create → running → expired → cleanup
3. **Concurrent Submission**: same flag submitted simultaneously, only counts once
4. **WebSocket**: multiple clients receive leaderboard updates simultaneously
5. **RBAC**: Player cannot access Admin endpoints

---

## Git Workflow

### Branch Strategy

```
main         Production branch, protected
 ↑
develop      Development branch, integration testing
 ↑
feature/*    Feature branches (e.g., feature/flag-submit)
 ↑
hotfix/*     Emergency fixes
```

### Merge Rules

1. **No push to main**: All changes via PR
2. **Squash Merge**: Feature branches squash into develop
3. **Rebase**: develop rebased onto main for linear history
4. **Code Review**: At least 1 approval required before merging

### PR Template

```markdown
## Description
<!-- What was done and why -->

## Testing
<!-- How to verify -->

## Checklist
- [ ] TDD: failing test written first
- [ ] All tests pass
- [ ] No lint errors
- [ ] Documentation updated
- [ ] No sensitive information leaked
```

---

## Release Process

### Versioning

Follows SemVer: `MAJOR.MINOR.PATCH`

- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

### Release Steps

1. Create `release/vX.Y.Z` from `develop`
2. Update `CHANGELOG.md`
3. PR → `main`, tag after merge
4. GitHub Actions automatically builds and pushes to GHCR
5. Publish GitHub Release

### Docker Image

```bash
# Build
docker build -t ghcr.io/gandli/yctf:latest .

# Run
docker compose up -d
```

Image tags:
- `latest` — main branch
- `vX.Y.Z` — release version
- `dev` — develop branch

---

## Troubleshooting

### Common Issues

| Symptom | Troubleshooting |
|---------|-----------------|
| Backend fails to start | Check PG/Redis connection, port conflicts |
| Container won't start | `docker ps -a` + `docker logs` |
| Frontend HMR broken | Check Vite proxy configuration |
| Test timeout | `go test -timeout 60s` |

### Debugging

```bash
# Backend Delve debugging
dlv debug cmd/server/main.go

# Frontend React DevTools
# Browser extension
```
