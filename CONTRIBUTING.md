# Contributing to YCTF

Welcome and thank you for considering contributing to YCTF!

## Code of Conduct

Be respectful, constructive, and inclusive. Harassment or hate speech will not be tolerated.

## How to Contribute

### Reporting Bugs

1. Check if the issue already exists
2. Open a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - OS and version info

### Suggesting Features

1. Open a Discussion first for major features
2. Explain the use case and expected behavior

### Pull Requests

1. Fork the repo and create a branch (`feature/my-feature`)
2. Follow TDD: failing test → code → passing test
3. Ensure all tests pass
4. Update documentation if needed
5. Submit PR with clear description

## Development Setup

### Backend (Go)

```bash
cd src
go mod tidy
go run cmd/server/main.go
```

### Frontend (React)

```bash
cd clientapp
npm install
npm run dev
```

### Full Stack (Docker)

```bash
docker compose up -d
```

## Style Guidelines

- **Go**: Follow standard Go conventions (`gofmt`, `golangci-lint`)
- **React**: Functional components, hooks, TypeScript strict mode
- **Commits**: Conventional commits (`feat:`, `fix:`, `docs:`, `chore:`)

## License

By contributing, you agree your contributions will be licensed under AGPL-3.0.
