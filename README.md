# YCTF 🚬🏁

**YCTF** — A modern, open-source CTF (Capture The Flag) competition platform.

> "烟起旗扬" — Where smoke rises, the flag flies.

---

## 🎯 What is YCTF?

YCTF is a full-featured Jeopardy-style CTF competition platform built with Go and React. It supports dynamic container deployment, real-time scoreboards, dynamic flag injection, and multi-language support.

Designed for security teams, CTF organizers, and training environments who want to host professional-grade competitions with minimal infrastructure overhead.

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| **Backend** | Go 1.25 (chi router, pgx, go-redis, gorilla/websocket) |
| **Frontend** | React 19 + Vite + Mantine + Tailwind CSS |
| **Database** | PostgreSQL |
| **Cache/Realtime** | Redis (Sorted Sets for leaderboard, rate limiting) |
| **Runtime** | Docker (dynamic container distribution) |
| **Deploy** | docker-compose + GitHub Actions → GHCR |

---

## ✨ Key Features

- 🏁 **Dynamic Container Distribution** — Each team gets isolated challenge instances with unique flags injected via environment variables
- 📊 **Real-time Leaderboard** — WebSocket-powered live scoreboard with Redis Sorted Sets
- 📉 **Dynamic Scoring** — Scores decrease as more teams solve challenges (first-blood bonus)
- 🌐 **i18n** — Full Chinese/English internationalization
- 🔒 **Three-tier RBAC** — Admin, Author, Player with optional team review
- 📝 **Writeup Collection** — Configurable post-competition writeup submission
- ⚡ **One-command Deploy** — `docker compose up -d` and you're running

---

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [PRD.md](docs/PRD.md) | Product Requirements Document |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | System Architecture |
| [DECISIONS.md](docs/DECISIONS.md) | Technical Decision Records (ADR) |
| [API.md](docs/API.md) | REST API & WebSocket Documentation |
| [DATABASE.md](docs/DATABASE.md) | Database Schema & Redis Structures |
| [DEVELOPMENT.md](docs/DEVELOPMENT.md) | Development Setup & Coding Standards |
| [SECURITY.md](docs/SECURITY.md) | Security Design & Hardening |
| [TODO.md](docs/TODO.md) | Development Roadmap |
| [CONTRIBUTING.md](docs/CONTRIBUTING.md) | Contributing Guidelines |
| [LICENSE](docs/LICENSE) | AGPL-3.0 License |

---

## 🚀 Quick Start

### Prerequisites

- Docker 24.0+
- Docker Compose v2

### Launch

```bash
# Clone
git clone https://github.com/gandli/yctf.git && cd yctf

# Start all services
docker compose up -d

# Access
# Frontend (dev server): http://localhost:5173
# Backend API: http://localhost:8080
```

---

## 🤝 Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md) for guidelines.

---

## 📄 License

AGPL-3.0 — See [LICENSE](docs/LICENSE).

---

## 🌟 Acknowledgments

Inspired by and learned from:
- [GZCTF](https://github.com/GZTimeWalker/GZCTF) — Feature design reference
- [rCTF](https://github.com/otter-sec/rctf) — Architecture simplicity
- [A1CTF](https://github.com/carbofish/A1CTF) — Project structure reference
