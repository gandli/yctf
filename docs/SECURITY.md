# YCTF Security Documentation

> Security design, threat model, and hardening measures.

---

## Threat Model

### Attacker Profiles

| Attacker | Goal | Methods |
|----------|------|---------|
| Cheating player | Share flags, brute force | Submit others' flags, high-frequency attempts |
| External attacker | Obtain flags, disrupt competition | Container escape, API brute force |
| Malicious insider | Tamper scores, leak challenges | Privilege escalation, data export |
| Script kiddie | Denial of service | CC attacks, slowloris |

### Critical Assets

| Asset | Protection Level | Description |
|-------|------------------|-------------|
| Flag signing key | 🔴 High | Leak allows forging any flag |
| Database | 🔴 High | User credentials, flag records |
| Admin credentials | 🔴 High | Platform control |
| Challenge source | 🟡 Medium | Early leak affects fairness |
| Frontend code | 🟢 Low | Public repo, no sensitive data |

---

## Security Measures

### Flag Security

| Measure | Implementation |
|---------|---------------|
| Uniqueness | Per-team per-challenge unique flag, HMAC(team_id + challenge_id + secret) |
| Anti-forgery | HMAC-SHA256, key exists only server-side |
| Anti-replay | Correct flag counts once, duplicate submissions ignored |
| Anti-enumeration | Format validation + length limit + rate limiting |
| Injection defense | Parameterized queries, flags never concatenated into SQL |

```go
// Flag generation example
func GenerateFlag(teamID, challengeID, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(teamID + ":" + challengeID))
    hash := hex.EncodeToString(h.Sum(nil))[:16]
    return fmt.Sprintf("flag{%s}", hash)
}
```

### Container Security

| Measure | Implementation |
|---------|---------------|
| Least privilege | `--cap-drop=ALL`, only add necessary capabilities |
| Read-only rootfs | `--read-only` |
| No privilege escalation | `--security-opt=no-new-privileges:true` |
| Resource limits | CPU / Memory / PidsLimit |
| Network isolation | Independent Docker network, no external access |
| Log audit | Container logs to stdout/stderr |

```go
// Secure container creation config
containerConfig := &container.Config{
    Image: image,
    Env:   []string{fmt.Sprintf("FLAG=%s", flag)},
}

hostConfig := &container.HostConfig{
    CapDrop:          []string{"ALL"},
    ReadonlyRootfs:   true,
    SecurityOpt:      []string{"no-new-privileges:true"},
    Resources: container.Resources{
        NanoCPUs: 500000000,    // 0.5 CPU
        Memory:   128 * 1024 * 1024,  // 128MB
        PidsLimit: 50,
    },
}
```

### API Security

| Measure | Implementation |
|---------|---------------|
| Authentication | JWT + Bearer Token |
| Rate limiting | Redis sliding window (10 req/min per user) |
| CORS | Whitelist origins |
| SQL injection | pgx parameterized queries |
| XSS | React auto-escape + CSP headers |
| CSRF | SameSite cookie + Origin validation |
| Input validation | chi middleware + validator library |

### Authentication Security

| Measure | Implementation |
|---------|---------------|
| Password storage | bcrypt (cost=12) |
| Token expiration | Access 15min / Refresh 7d |
| Refresh rotation | Each refresh generates new refresh token |
| Token revocation | Redis blacklist (added on logout) |

---

## Incident Response

### Severity Levels

| Level | Incident | Response |
|-------|----------|----------|
| P0 | Flag key leaked | Rotate key immediately, regenerate all flags |
| P0 | Container escape | Isolate affected containers, notify players |
| P1 | User data breach | Force password reset, notify affected users |
| P1 | DDoS attack | Enable CDN/WAF, rate limiting |
| P2 | Abnormal submission pattern | Review logs, ban if necessary |

### Contact

Security incidents: email `security@chenxuexin.com`

---

## Security Audit Checklist

- [ ] Flag key rotated (default → custom)
- [ ] All API endpoints require authentication
- [ ] Container security config verified
- [ ] Database not exposed to public internet
- [ ] Redis password set
- [ ] HTTPS configured
- [ ] CORS whitelist restricted
- [ ] Logs exclude sensitive data (flags, passwords)
