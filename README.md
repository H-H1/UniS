# uniS — 微信小程序后端服务

基于 Go + Gin + GORM + MySQL 的微信小程序后端，提供微信登录（code2session）、JWT 鉴权等核心功能，日志按天写入本地文件。

---

## 技术栈

| 层级 | 技术 |
|------|------|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| ORM | [GORM](https://gorm.io) v1.31 + MySQL 驱动 |
| 认证 | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) v5 |
| 微信 | 自封装 `pkg/wechat`，调用 `jscode2session` 接口 |
| 日志 | 自封装 `pkg/logger`，按天滚动写入 `logs/` 目录 |
| 语言版本 | Go 1.25+ |

---

## 项目结构

```
.
├── cmd/server/          # 程序入口 main.go
├── config/              # 配置加载（环境变量优先）
├── datascripts/         # 数据库初始化 SQL
├── internal/
│   ├── database/        # GORM 连接与 AutoMigrate
│   ├── handler/         # HTTP 处理层（请求解析、响应）
│   ├── middleware/       # JWT 鉴权中间件
│   ├── model/           # 数据模型（User）
│   ├── repository/      # 数据访问层（CRUD）
│   ├── router/          # 路由注册
│   └── service/         # 业务逻辑层
├── pkg/
│   ├── jwt/             # JWT 签发与解析
│   ├── logger/          # 结构化日志（文件 + 标准输出）
│   └── wechat/          # 微信 API 客户端
├── logs/                # 运行时日志（自动生成，已 gitignore）
├── .env.example         # 环境变量示例
└── go.mod
```

---

## 快速启动

### 1. 克隆并安装依赖

```bash
git clone https://github.com/your-org/uniS.git
cd uniS
go mod tidy
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env，填入真实的数据库连接和微信 AppSecret
```

| 变量 | 说明 | 示例 |
|------|------|------|
| `SERVER_PORT` | 监听端口 | `8099` |
| `MYSQL_USERNAME` | 数据库用户名 | `root` |
| `MYSQL_PASSWORD` | 数据库密码 | `your_password` |
| `MYSQL_ADDRESS` | 数据库地址 | `127.0.0.1:3306` |
| `MYSQL_DATABASE` | 数据库名 | `unis` |
| `WX_APP_ID` | 微信小程序 AppID | `wx...` |
| `WX_APP_SECRET` | 微信小程序 AppSecret | `...` |
| `JWT_SECRET` | JWT 签名密钥 | 随机长字符串 |

### 3. 初始化数据库

```bash
mysql -u root -p < datascripts/init.sql
```

> 也可以跳过此步骤，GORM AutoMigrate 会在服务启动时自动建表。

### 4. 启动服务

```bash
go run cmd/server/main.go
```

服务默认监听 `:8099`。

---

## API 接口

### 公开接口

#### 微信登录

```
POST /auth/wx-login
```

**请求体**

```json
{
  "code":      "微信 wx.login() 返回的 code（必填）",
  "nickName":  "用户昵称",
  "avatarUrl": "头像地址",
  "gender":    1,
  "country":   "China",
  "province":  "广东",
  "city":      "深圳"
}
```

**成功响应**

```json
{
  "code": 0,
  "msg":  "ok",
  "data": {
    "token": "eyJhbGci...",
    "user": {
      "id":         42,
      "open_id":    "oXXXX",
      "nick_name":  "张三",
      "avatar_url": "https://...",
      "gender":     1,
      "country":    "China",
      "province":   "广东",
      "city":       "深圳",
      "created_at": "2026-05-01T17:48:40Z",
      "updated_at": "2026-05-01T17:48:40Z"
    }
  }
}
```

### 鉴权接口

所有 `/api/*` 路由需在请求头携带 JWT：

```
Authorization: Bearer <token>
```

#### 获取当前用户信息

```
GET /api/me
```

**响应**

```json
{
  "user_id": 42,
  "open_id": "oXXXX"
}
```

---

## 日志

日志同时输出到标准输出和 `logs/app-YYYY-MM-DD.log`，按天自动滚动。

```
[2026-05-01 17:48:40.746] [INFO ] [auth_handler] 收到 WxLogin 请求 | client_ip=192.168.5.7 nick_name=张三
[2026-05-01 17:48:40.977] [ERROR] [wechat] HTTP 请求失败 | appid=wx... error=...
```

日志级别：`INFO` / `WARN` / `ERROR`

---

## 部署

项目根目录包含 `Dockerfile`，支持容器化部署：

```bash
docker build -t unis-server .
docker run -p 8099:8099 --env-file .env unis-server
```

---

## 开发规范

详见 [go_spec.MD](go_spec.MD)（后端）和 [uniapp_spec.MD](uniapp_spec.MD)（前端）。

---

## License

[MIT](LICENSE)
