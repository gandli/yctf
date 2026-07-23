# YCTF API 文档

> RESTful API 设计与 WebSocket 事件规范。

---

## 基础信息

| 项目 | 值 |
|------|-----|
| Base URL | `/api/v1` |
| 认证 | Bearer Token (JWT) |
| 格式 | JSON |
| CORS | 配置化（默认 `http://localhost:5173`） |

---

## 认证

### POST `/auth/register`

注册新用户。

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

## 用户

### GET `/users/me`

获取当前用户信息。

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

更新用户信息（昵称、密码）。

---

## 团队

### POST `/teams`

创建团队。

**Request:**
```json
{
  "name": "TeamA"
}
```

### POST `/teams/join`

通过邀请码加入团队。

**Request:**
```json
{
  "invite_code": "XXXXXX"
}
```

### GET `/teams/:id`

获取团队详情（成员、解题历史）。

---

## 题目

### GET `/challenges`

获取题目列表（按分类筛选）。

**Query Params:**
- `category` — web/pwn/crypto/re/misc/osint
- `page` — 页码
- `limit` — 每页数量

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

获取题目详情。

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

创建题目。

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

## 容器实例

### POST `/challenges/:id/start`

启动容器实例。

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

停止容器实例。

### GET `/instances/:id`

获取实例状态。

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

## Flag 提交

### POST `/submit`

提交 flag。

**Request:**
```json
{
  "challenge_id": "uuid",
  "flag": "flag{sql_injection_success}"
}
```

**Response (200) — 正确:**
```json
{
  "correct": true,
  "message": "🎉 Flag correct! +100 points",
  "score_gained": 100,
  "team_score": 1600
}
```

**Response (200) — 错误:**
```json
{
  "correct": false,
  "message": "❌ Incorrect flag"
}
```

**Response (429) — 限速:**
```json
{
  "error": "rate_limit_exceeded",
  "message": "Too many submissions. Try again in 30 seconds.",
  "retry_after": 30
}
```

---

## 排行榜

### GET `/scoreboard`

获取排行榜。

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

获取团队得分时间线（用于图表）。

**Query Params:**
- `team_id` — 可选，筛选特定团队

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

## Writeup

### POST `/writeups`

提交 writeup。

**Request:**
```json
{
  "challenge_id": "uuid",
  "url": "https://example.com/writeup",
  "content": "optional markdown"
}
```

### GET `/writeups`

获取已审核的 writeup 列表。

---

## 管理端点 (Admin)

### GET `/admin/users`

用户列表（分页、搜索）。

### PATCH `/admin/users/:id`

更新用户角色/封禁。

### GET `/admin/submissions`

提交记录审计。

### GET `/admin/stats`

平台统计（用户数、提交数、容器数）。

---

## WebSocket 事件

**连接:** `ws://localhost:8080/ws`

**握手:** 携带 JWT token（query param 或 header）

### 事件类型

| 事件 | 方向 | Payload |
|------|------|---------|
| `score_update` | Server → Client | `{ teams: [...], timestamp }` |
| `challenge_solved` | Server → Client | `{ team_name, challenge_title, points }` |
| `notification` | Server → Client | `{ type, message }` |
| `ping` | Bidirectional | `{}` |

### 示例

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

## 错误格式

```json
{
  "error": "error_code",
  "message": "Human readable message",
  "details": {}
}
```

### 错误码

| HTTP | Code | 说明 |
|------|------|------|
| 400 | `bad_request` | 请求格式错误 |
| 401 | `unauthorized` | 未认证 |
| 403 | `forbidden` | 权限不足 |
| 404 | `not_found` | 资源不存在 |
| 409 | `conflict` | 资源冲突（如重复注册） |
| 429 | `rate_limit_exceeded` | 请求过于频繁 |
| 500 | `internal_error` | 服务器内部错误 |
