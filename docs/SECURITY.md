# YCTF 安全文档

> 安全设计、威胁模型与加固措施。

---

## 威胁模型

### 攻击者画像

| 攻击者 | 目标 | 手段 |
|--------|------|------|
| 参赛作弊者 | 共享 flag、暴力破解 | 提交他人 flag、高频尝试 |
| 外部攻击者 | 获取 flag、破坏比赛 | 容器逃逸、API 爆破 |
| 内部恶意用户 | 篡改成绩、泄露题目 | 权限提升、数据导出 |
| 脚本小子 | 拒绝服务 | CC 攻击、慢连接 |

### 关键资产

| 资产 | 保护级别 | 说明 |
|------|----------|------|
| Flag 生成密钥 | 🔴 高 | 一旦泄漏可伪造任意 flag |
| 数据库 | 🔴 高 | 用户凭证、flag 记录 |
| 管理员凭证 | 🔴 高 | 平台控制权 |
| 挑战源码 | 🟡 中 | 提前泄漏影响公平 |
| 前端代码 | 🟢 低 | 公开仓库，无敏感信息 |

---

## 安全措施

### Flag 安全

| 措施 | 实现 |
|------|------|
| 唯一性 | 每队每题唯一 flag，HMAC(team_id + challenge_id + secret) |
| 防伪造 | HMAC-SHA256，密钥仅存在于服务端 |
| 防重放 | 正确 flag 仅记一次分，重复提交不计 |
| 防遍历 | 格式校验 + 长度限制 + 速率限制 |
| 注入防御 | 参数化查询，flag 不参与 SQL 拼接 |

```go
// flag 生成示例
func GenerateFlag(teamID, challengeID, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(teamID + ":" + challengeID))
    hash := hex.EncodeToString(h.Sum(nil))[:16]
    return fmt.Sprintf("flag{%s}", hash)
}
```

### 容器安全

| 措施 | 实现 |
|------|------|
| 权限最小化 | `--cap-drop=ALL`，仅保留必要能力 |
| 只读根文件系统 | `--read-only` |
| 禁止提权 | `--security-opt=no-new-privileges:true` |
| 资源限制 | CPU / Memory / PidsLimit |
| 网络隔离 | 独立 Docker network，禁止外网 |
| 日志审计 | 容器日志发送到 stdout/stderr |

```go
// 容器创建安全配置
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

### API 安全

| 措施 | 实现 |
|------|------|
| 认证 | JWT + Bearer Token |
| 速率限制 | Redis 滑动窗口（10 req/min per user） |
| CORS | 白名单 origin |
| SQL 注入 | pgx 参数化查询 |
| XSS | React 自动转义 + CSP 头 |
| CSRF | SameSite cookie + 验证 Origin |
| 输入验证 | chi middleware + validator 库 |

### 认证安全

| 措施 | 实现 |
|------|------|
| 密码存储 | bcrypt (cost=12) |
| Token 过期 | Access 15min / Refresh 7d |
| 刷新轮换 | 每次 refresh 生成新的 refresh token |
| Token 撤销 | Redis 黑名单（登出时加入） |

---

## 安全响应

### 事件分级

| 级别 | 事件 | 响应 |
|------|------|------|
| P0 | Flag 密钥泄漏 | 立即轮换密钥、重新生成所有 flag |
| P0 | 容器逃逸 | 隔离受影响容器、通知玩家 |
| P1 | 用户数据泄露 | 强制密码重置、通知受影响用户 |
| P1 | DDoS 攻击 | 启用 CDN/WAF、限流 |
| P2 | 异常提交模式 | 审查日志、必要时封禁 |

### 联系

安全事件请发送邮件至 `security@chenxuexin.com`。

---

## 安全审计清单

- [ ] Flag 密钥已轮换（默认值 → 自定义）
- [ ] 所有 API 端点需要认证
- [ ] 容器安全配置已验证
- [ ] 数据库无公网暴露
- [ ] Redis 设置了密码
- [ ] HTTPS 已配置
- [ ] CORS 白名单已限制
- [ ] 日志不记录敏感信息（flag、密码）
