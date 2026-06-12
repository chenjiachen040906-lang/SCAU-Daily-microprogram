# SCAU Daily

华南农业大学学生的"校园大脑外挂" —— 课表查询、智能待办、校园信息聚合。

## 项目结构

```
├── server/          # Go 后端 (Fiber + PostgreSQL + Redis)
├── miniapp/         # 微信小程序 (uni-app, 待开发)
└── docker-compose.yml
```

## 快速启动

### 前置条件

- Docker & Docker Compose
- Go 1.21+（本地开发时需要）

### 使用 Docker Compose 启动

```bash
docker compose up -d
```

服务将在 `http://localhost:8080` 启动。

### 本地开发

```bash
cd server
cp .env.example .env
# 确保 PostgreSQL 和 Redis 已运行
go run ./cmd/server
```

## API 文档

启动后访问 `GET /health` 确认服务状态。

所有 API 路径前缀为 `/api/v1`，需要 Bearer Token 认证（登录接口除外）。

## 技术栈

- **后端**: Go + Fiber + GORM + PostgreSQL + Redis
- **小程序**: uni-app (Vue 3 + TypeScript)（待开发）
- **App**: Flutter 3.x（待开发）
