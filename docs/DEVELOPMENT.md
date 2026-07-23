# YCTF 开发者文档

> 开发环境搭建、代码规范、测试策略、部署指南。

---

## 目录

- [开发环境 Setup](#开发环境-setup)
- [代码规范](#代码规范)
- [测试策略](#测试策略)
- [Git 工作流](#git-工作流)
- [发布流程](#发布流程)

---

## 开发环境 Setup

### 前置依赖

| 工具 | 版本 | 用途 |
|------|------|------|
| Go | 1.24+ | 后端开发 |
| Node.js | 20+ | 前端开发 |
| Docker | 24+ | 容器运行时 |
| Docker Compose | v2+ | 本地开发环境 |
| Make | 任意 | 任务编排 |

### 快速开始

```bash
# 克隆
git clone https://github.com/gandli/yctf.git && cd yctf

# 启动依赖服务（PG + Redis）
docker compose -f compose.dev.yml up -d

# 后端
cd src
go mod tidy
go run cmd/server/main.go

# 前端（新终端）
cd clientapp
npm install
npm run dev
```

### 访问

| 服务 | URL |
|------|-----|
| 前端 Dev Server | http://localhost:5173 |
| 后端 API | http://localhost:8080 |
| Grafana | http://localhost:3001 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

### 环境变量

```bash
# 后端 (.env)
DATABASE_URL=postgres://yctf:yctf@localhost:5432/yctf?sslmode=disable
REDIS_URL=redis://localhost:6379/0
JWT_SECRET=dev-secret-change-in-production
PORT=8080
CORS_ORIGINS=http://localhost:5173

# 前端 (.env)
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

---

## 代码规范

### Go 规范

| 规则 | 说明 |
|------|------|
| 格式化 | `gofmt -w`（CI 强制） |
| Lint | `golangci-lint run`（配置 `.golangci.yml`） |
| 命名 | PascalCase 导出、camelCase 未导出、全大写缩写（HTTP、ID） |
| 错误处理 | `if err != nil` 必须处理，禁止 `_ = err` |
| 接口 | 单一职责，文件名 `interface.go` |
| 测试 | 同目录 `*_test.go`，表驱动测试优先 |
| 依赖注入 | 通过接口，禁止全局变量（除 `const`） |

### React/TypeScript 规范

| 规则 | 说明 |
|------|------|
| 组件 | 函数组件 + Hooks，禁止 class 组件 |
| 文件命名 | PascalCase 组件、camelCase 工具 |
| 类型 | 严格模式，禁止 `any`（除第三方无类型） |
| Props 接口 | `interface ComponentNameProps` |
| 样式 | Tailwind 优先，Mantine `sx` 次之 |
| 导入 | 绝对路径 `@/` → `src/` |
| 状态 | Zustand stores，避免 prop drilling |

### Commit 规范

```
feat:     新功能
fix:      修复 Bug
docs:     文档变更
style:    代码格式（不影响功能）
refactor: 重构
perf:     性能优化
test:     测试相关
chore:    构建/依赖/CI 配置
ci:       CI 配置
```

示例：`feat: add flag submission rate limiting`

---

## 测试策略

### TDD 铁律

> **没有失败的测试先行，就没有生产代码。**

```bash
# 后端测试
cd src
go test ./... -v -count=1              # 全量
go test ./controllers/ -run TestSubmit  # 指定

# 前端测试
cd clientapp
npm run test                            # Vitest 单元
npm run test:e2e                        # Playwright E2E
```

### 测试金字塔

```
        /  E2E (Playwright)  \         ← 用户旅程
       / Integration (API Test) \       ← 模块协作
      /    Unit (Go test / Vitest) \    ← 基础单元
```

### 覆盖率要求

| 模块 | 目标 |
|------|------|
| `utils/` | 90%+ |
| `controllers/` | 80%+ |
| `db/queries/` | 85%+ |
| `middleware/` | 75%+ |
| 前端 hooks | 80%+ |

### 关键测试场景

1. **Flag 提交**：正确 flag → 得分；错误 flag → 拒收；限速 → 429
2. **容器生命周期**：创建 → 运行 → 过期 → 清理
3. **并发提交**：同一 flag 同时提交，只记一次
4. **WebSocket**：多客户端同时收到排行榜更新
5. **RBAC**：Player 不能访问 Admin 端点

---

## Git 工作流

### 分支策略

```
main         生产分支，受保护
  ↑
develop      开发分支，集成测试
  ↑
feature/*    功能分支（例：feature/flag-submit）
  ↑
hotfix/*     紧急修复
```

### 合并规则

1. **禁推 main**：所有变更通过 PR
2. **Squash Merge**：功能分支合入 develop 时压缩
3. **Rebase**：develop 合入 main 时保持线性
4. **Code Review**：至少 1 人 approve 才能合并

### PR 模板

```markdown
## 改动描述
<!-- 做了什么，为什么 -->

## 测试方式
<!-- 如何验证 -->

## Checklist
- [ ] TDD：先写失败测试
- [ ] 所有测试通过
- [ ] 无 lint 错误
- [ ] 文档已更新
- [ ] 无敏感信息泄漏
```

---

## 发布流程

### 版本号

遵循 SemVer：`MAJOR.MINOR.PATCH`

- MAJOR：不兼容变更
- MINOR：功能新增（向下兼容）
- PATCH：Bug 修复

### 发布步骤

1. 从 `develop` 创建 `release/vX.Y.Z`
2. 更新 `CHANGELOG.md`
3. PR → `main`，合并后打 tag
4. GitHub Actions 自动构建镜像并推送到 GHCR
5. 发布 GitHub Release

### Docker 镜像

```bash
# 构建
docker build -t ghcr.io/gandli/yctf:latest .

# 运行
docker compose up -d
```

镜像标签：
- `latest` — main 分支
- `vX.Y.Z` — 发布版本
- `dev` — develop 分支

---

## 故障排查

### 常见问题

| 现象 | 排查 |
|------|------|
| 后端启动失败 | 检查 PG/Redis 连接、端口占用 |
| 容器无法启动 | `docker ps -a` + `docker logs` |
| 前端热重载失效 | 检查 Vite proxy 配置 |
| 测试超时 | `go test -timeout 60s` |

### 调试

```bash
# 后端 Delve 调试
dlv debug cmd/server/main.go

# 前端 React DevTools
# 浏览器扩展即可
```
