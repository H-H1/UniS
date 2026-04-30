

# Go 后端开发规范文档（Backend Spec）

## 1. 项目概述

- **项目名称**：____________（请填写）
- **技术栈**：Go 1.21+ / Gin / Echo / GORM / sqlc / MySQL/PostgreSQL / Redis / JWT
- **架构模式**：RESTful API + 分层架构（Handler → Service → Repository）
- **核心目标**：高性能、可维护、可测试

---

## 2. 项目结构

```
.
├── cmd/
│   └── server/            # 程序入口（main.go）
├── internal/              # 私有代码
│   ├── handler/           # HTTP 处理层
│   ├── service/           # 业务逻辑层
│   ├── repository/        # 数据访问层
│   ├── model/             # 数据模型
│   ├── dto/               # 请求/响应结构体
│   └── middleware/        # 鉴权、日志、限流等
├── pkg/                   # 可复用工具库
├── config/                # 配置加载
├── migrations/            # SQL 迁移文件
├── scripts/               # 构建/部署脚本
├── test/                  # 集成测试
├── docs/                  # API 文档
├── go.mod / go.sum
└── Makefile
```

---

## 3. 编码规范

- 必须通过：`go fmt`、`go vet`、`staticcheck`、`golangci-lint`
- 命名：
  - 文件：小写+下划线（`user_service.go`）
  - 包：小写单数（`handler`）
  - 公开标识：大驼峰（`UserService`）
  - 私有标识：小驼峰（`getUserByID`）
- 错误处理：显式处理 error，禁止忽略关键错误
- 避免 panic：仅限不可恢复的初始化阶段
- 注释：所有导出函数/结构体/常量必须有注释
- import 顺序：标准库 → 第三方 → 本项目包（空行分组）

---

## 4. API 设计规范（RESTful）

- 路径：`/v1/resources`（名词复数）
- 方法：
  - `GET /resources`
  - `GET /resources/{id}`
  - `POST /resources`
  - `PUT|PATCH /resources/{id}`
  - `DELETE /resources/{id}`
- 请求/响应：JSON，字段小驼峰（`userName`）
- 状态码：200 / 201 / 400 / 401 / 403 / 404 / 500
- 错误响应格式：
```json
{
  "code": 10001,
  "message": "用户名已存在",
  "details": null
}
```
- 文档：Swagger/OpenAPI（推荐 `swaggo/swag`）

---

## 5. 分层职责

| 层       | 职责                         | 允许依赖               |
|----------|------------------------------|------------------------|
| Handler  | 请求解析、校验、调用 Service | Service                |
| Service  | 业务逻辑、事务、调用 Repo    | Repository、其他 Service |
| Repository | 数据库 CRUD                | Model                  |
| Model    | 实体/表结构定义              | 无                     |

> Handler 不写业务逻辑，Service 不写 SQL，Repository 不写业务判断。

---

## 6. 数据库与迁移

- 连接池：`database/sql` + `MaxOpenConns` / `MaxIdleConns`
- 迁移工具：`golang-migrate` / `atlas`
- SQL 方式：
  - 推荐 `sqlc` 生成类型安全代码
  - 使用 GORM 时避免 N+1，显式 `Preload`
- 事务：在 Service 层控制，Repository 支持传入 `*sql.Tx` 或事务闭包

---

## 7. 配置管理

- 库：`viper`
- 优先级：环境变量 > 配置文件
- 敏感信息（密钥、密码）仅通过环境变量注入，不提交到仓库
- 示例：
  - `config/config.yaml`
  - `config/config.local.yaml`（gitignored）

---

## 8. 日志与可观测性

- 推荐：`slog` / `zap` / `zerolog`
- 级别：Debug / Info / Warn / Error
- 格式：JSON（结构化字段，如 `user_id`）
- 每个请求：唯一 `request_id`（中间件生成）

---

## 9. 安全规范

- 认证：JWT（Access + Refresh）
- 授权：RBAC
- SQL 注入：参数化查询，禁止拼接 SQL
- 限流：令牌桶 / 漏桶（本地或 Redis）
- HTTP 超时：`ReadTimeout` / `WriteTimeout` / `IdleTimeout`
- 密码：bcrypt 哈希，禁止返回前端

---

## 10. 测试规范

- 单元测试：`testing` + `testify/assert` + `mockery`
- 覆盖率：核心业务 ≥ 80%
- 集成测试：`testcontainers-go`（真实 DB / Redis）
- 命令：
```bash
go test -v -race -cover ./...
```

---

## 11. 性能与可靠性

- pprof：`import _ "net/http/pprof"`
- 优雅关闭：监听 `SIGINT`/`SIGTERM`，等待现有请求结束
- 健康检查：`/health`

---

## 12. 构建与部署

Makefile 示例目标：
```makefile
build:
	go build -o bin/server cmd/server/main.go

test:
	go test -race -cover ./...

migrate-up:
	migrate -path migrations -database "$DB_URL" up
```

- Docker：多阶段构建，最终镜像 `scratch` / `alpine`
- 环境：dev → staging → production

---

## 13. 依赖管理

- Go Modules
- 定期执行 `go mod tidy`
- 关键依赖锁定版本，避免 `@latest`
- 优先标准库或轻量库

---

## 14. 文档要求

- `README.md`：简介 / 快速启动 / 配置 / API 入口
- `docs/api/`：Swagger 或 Markdown
- 复杂模块：`internal/xxx/README.md`

---

## 15. 版本与变更管理

- API 路径版本：`/v1/...`、`/v2/...`
- 重大变更需兼容期或迁移指南
- Git 标签：语义化版本（`v1.2.3`）

---

**如果需要，我可以再帮你把这份 spec 输出为 Markdown 文件内容，方便你直接复制到 `.md` 或 `spec/` 目录中。**