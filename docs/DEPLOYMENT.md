# YCTF 部署指南

> 从源码到生产的完整部署流程。

---

## 快速部署（Docker Compose）

### 前置要求

- Docker 24.0+
- Docker Compose v2
- 2GB 可用内存
- 10GB 可用磁盘

### 一键启动

```bash
# 克隆
git clone https://github.com/gandli/yctf.git && cd yctf

# 启动
docker compose up -d

# 访问
# API: http://localhost:8080
# 健康检查: http://localhost:8080/health
```

### 验证部署

```bash
# 检查所有服务状态
docker compose ps

# 查看日志
docker compose logs -f server

# 测试 API
curl http://localhost:8080/health
```

---

## 生产部署

### 1. 环境变量配置

创建 `.env` 文件：

```bash
# 数据库密码（必须修改）
POSTGRES_PASSWORD=your_secure_password_here

# Redis 密码（必须修改）
REDIS_PASSWORD=your_secure_redis_password

# JWT 密钥（必须修改，至少 32 字节）
JWT_SECRET=your-very-long-secret-key-here-min-32-bytes

# Flag 签名密钥（必须修改）
FLAG_SECRET=another-secure-secret-for-flag-signing

# CORS 允许的域名
CORS_ORIGINS=https://your-domain.com
```

### 2. 启动生产环境

```bash
docker compose -f docker-compose.prod.yml up -d
```

### 3. 配置 Nginx（可选）

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }

    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
}
```

### 4. SSL 证书（Let's Encrypt）

```bash
sudo certbot --nginx -d your-domain.com
```

---

## 从源码构建

### 本地开发

```bash
# 后端
cd src
go mod tidy
go run cmd/server/main.go

# 前端（需要 Node.js 20+）
cd clientapp
npm install
npm run dev
```

### 构建 Docker 镜像

```bash
docker build -t yctf-server:latest ./src
```

---

## 故障排查

| 问题 | 排查 |
|------|------|
| 服务无法启动 | `docker compose logs server` |
| 数据库连接失败 | 检查 `POSTGRES_PASSWORD` 是否一致 |
| Redis 连接失败 | 检查 `REDIS_PASSWORD` 是否一致 |
| 容器无法启动 | `docker compose ps -a` + `docker logs <id>` |

---

## 安全加固清单

- [ ] 修改所有默认密码
- [ ] 使用强 JWT_SECRET（≥32 字节）
- [ ] 使用强 FLAG_SECRET
- [ ] 配置 CORS_ORIGINS 白名单
- [ ] 启用 HTTPS
- [ ] 配置防火墙（仅开放 80/443）
- [ ] 定期备份数据库
- [ ] 配置日志轮转
