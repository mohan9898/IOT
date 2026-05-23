# IoT 设备管理器

基于 Go + Vue 3 构建的轻量级 IoT 设备管理平台。支持设备注册、远程控制、MQTT 通信、实时状态监控，开箱即用，适合中小型物联网项目快速部署。

## 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端框架 | Go 1.21 + Gin | 高性能 HTTP 框架，RESTful API |
| 前端框架 | Vue 3 + Vite + Tailwind CSS | 响应式 UI，极速热更新 |
| 数据库 | SQLite (go-sqlite3) | 嵌入式数据库，零配置，无需额外服务 |
| MQTT | Eclipse Paho MQTT Go | 设备双向通信，支持 TLS/SSL 加密 |
| 认证 | JWT (JSON Web Token) | 无状态认证，支持 Token 过期 |
| 构建 | Docker Multi-stage | 单容器部署，镜像体积小 |
| CI/CD | GitHub Actions + GHCR | 自动构建并发布 Docker 镜像 |

## 功能特性

- **设备管理** — 设备注册、发现、CRUD 操作，支持自定义设备类型（智能灯、传感器、执行器、控制器、摄像头、恒温器、智能开关）
- **远程控制** — 通过 MQTT 向设备下发控制指令（ON/OFF/AUTO/SET_THRESHOLD 等），记录指令历史
- **实时监控** — MQTT 状态上报自动更新设备状态，WebSocket 实时推送前端
- **用户认证** — JWT 登录/注册，bcrypt 密码加密，Token 过期管理
- **安全防护** — 速率限制（100次/分钟）、安全响应头（CSP/X-Frame/XSS）、CORS 跨域控制、非 root 容器运行
- **数据持久化** — SQLite 嵌入式存储，支持自动定时备份（VACUUM INTO 安全备份）
- **健康检查** — 内置 `/health` 端点 + Docker HEALTHCHECK 自动检测
- **优雅降级** — MQTT 不可用时应用程序仍可正常运行

## 项目结构

```
IOT/
├── backend/                     # Go 后端
│   ├── config/config.go         # 配置管理（环境变量覆盖）
│   ├── internal/
│   │   ├── api/api.go           # REST API 路由和处理器
│   │   ├── auth/jwt.go          # JWT 认证管理
│   │   ├── db/
│   │   │   ├── sqlite.go        # SQLite 数据库操作
│   │   │   └── backup.go        # 自动备份管理器
│   │   ├── mqtt/client.go       # MQTT 客户端连接
│   │   └── metrics/             # Prometheus 监控指标
│   └── main.go                  # 服务入口
├── frontend/                    # Vue 3 前端
│   └── src/
│       ├── components/          # Vue 组件
│       │   ├── LoginPage.vue    # 登录页面
│       │   ├── DeviceList.vue   # 设备列表
│       │   ├── DeviceCard.vue   # 设备卡片
│       │   ├── AddDevice.vue    # 添加设备
│       │   ├── DeviceTypeManager.vue  # 设备类型管理
│       │   └── AccountSettings.vue    # 账户设置
│       ├── services/api.js      # API 请求封装
│       └── store/index.js       # Pinia 状态管理
├── dist/                        # 前端构建产物（已打包，无需 Node.js 环境）
├── .github/workflows/           # GitHub Actions CI/CD
│   └── docker-publish.yml       # 自动构建 Docker 镜像
├── Dockerfile                   # 多阶段构建
├── docker-compose.yml           # Docker Compose 部署配置
├── .env.example                 # 环境变量模板
├── deploy.sh                    # 一键部署脚本
└── .dockerignore
```

---

## 快速部署

### 前置要求

- [Docker](https://docs.docker.com/get-docker/)（20.10 或更高版本）
- [Docker Compose](https://docs.docker.com/compose/install/)（v2 或更高版本）

### 方式一：一键部署脚本（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/mohan9898/IOT.git
cd IOT

# 2. 运行部署脚本
chmod +x deploy.sh
./deploy.sh
```

脚本会自动完成以下步骤：
1. 检查 Docker 环境
2. 从 `.env.example` 创建 `.env` 配置文件（如不存在）
3. 从 GitHub Container Registry 拉取最新镜像
4. 启动容器服务

启动完成后访问：**http://localhost:6116**

### 方式二：手动 Docker Compose 部署

```bash
# 1. 克隆仓库
git clone https://github.com/mohan9898/IOT.git
cd IOT

# 2. 创建并编辑环境变量
cp .env.example .env
# 编辑 .env，必须设置 JWT_SECRET 和 MQTT 凭证

# 3. 拉取镜像并启动
docker compose pull
docker compose up -d

# 4. 查看运行状态
docker compose ps

# 5. 查看日志
docker compose logs -f
```

启动完成后访问：**http://localhost:6116**

### 方式三：从源码构建

```bash
# 1. 克隆仓库
git clone https://github.com/mohan9898/IOT.git
cd IOT

# 2. 构建 Docker 镜像
docker build -t iot-manager .

# 3. 创建数据目录
mkdir -p data

# 4. 运行容器
docker run -d \
  --name iot-manager \
  -p 6116:6116 \
  -v $(pwd)/data:/app/data \
  -e JWT_SECRET="$(openssl rand -base64 32)" \
  -e MQTT_USERNAME="your-mqtt-username" \
  -e MQTT_PASSWORD="your-mqtt-password" \
  iot-manager
```

---

## 配置说明

所有配置通过环境变量设置。将 `.env.example` 复制为 `.env` 并修改对应值。

### 必填配置

| 环境变量 | 说明 | 示例 |
|----------|------|------|
| `JWT_SECRET` | JWT 签名密钥，**生产环境必须修改** | `openssl rand -base64 32` 生成 |

### MQTT 配置

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `MQTT_BROKER` | MQTT 服务器地址 | `d11aab19.ala.cn-hangzhou.emqxsl.cn` |
| `MQTT_PORT` | MQTT 端口 | `8883`（TLS）/ `1883`（TCP） |
| `MQTT_PROTOCOL` | 连接协议 | `ssl`（也可用 `tcp`、`ws`、`wss`） |
| `MQTT_USERNAME` | MQTT 用户名 | （空） |
| `MQTT_PASSWORD` | MQTT 密码 | （空） |
| `MQTT_TLS_ENABLED` | 是否启用 TLS | `true` |

### 数据库配置

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `DB_PATH` | 数据库文件路径 | `/app/data/iot.db` |
| `DB_BACKUP_ENABLE` | 是否启用自动备份 | `true` |
| `DB_BACKUP_PATH` | 备份文件路径 | `/app/data/backup` |
| `DB_BACKUP_HOURS` | 备份间隔（小时） | `24` |

### 服务器配置

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `SERVER_HOST` | 监听地址 | `0.0.0.0` |
| `SERVER_PORT` | 监听端口 | `6116` |
| `GIN_MODE` | Gin 运行模式 | `release`（也可设置 `debug`） |
| `CORS_ALLOWED_ORIGINS` | 跨域允许的源，逗号分隔 | （空，使用同源默认策略） |

### JWT 配置

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `JWT_SECRET` | JWT 签名密钥（**必填**） | （无默认值，必须设置） |
| `JWT_EXPIRES_HOURS` | Token 过期时间（小时） | `24` |

---

## API 接口文档

服务默认运行在 `http://localhost:6116`。

### 公开接口（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/health` | 健康检查，返回 `{"status":"ok"}` |
| `GET` | `/metrics` | Prometheus 监控指标 |
| `POST` | `/api/auth/register` | 用户注册 |
| `POST` | `/api/auth/login` | 用户登录 |

### 认证接口（需要 Bearer Token）

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/auth/user-info` | 获取用户信息 |
| `POST` | `/api/auth/update-account` | 修改用户名 |
| `POST` | `/api/auth/update-password` | 修改密码 |

### 设备接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/devices` | 获取所有设备列表 |
| `GET` | `/api/devices/:id` | 获取单个设备详情 |
| `POST` | `/api/devices` | 创建设备 |
| `PUT` | `/api/devices/:id` | 更新设备信息 |
| `DELETE` | `/api/devices/:id` | 删除设备 |
| `GET` | `/api/devices/stats` | 设备统计（在线/离线/总数） |

### 设备类型接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/device-types` | 获取所有设备类型 |
| `POST` | `/api/device-types` | 添加自定义设备类型 |
| `GET` | `/api/device-types/:type_id` | 获取指定设备类型 |

### 控制接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/control/send` | 向设备发送控制指令 |
| `POST` | `/api/control/threshold` | 设置智能灯亮度阈值 |
| `GET` | `/api/control/history/:device_id` | 获取设备指令历史 |

### WebSocket

| 路径 | 说明 |
|------|------|
| `GET` | `/api/ws` | 实时数据推送（需携带 Token） |

### API 使用示例

```bash
# 注册新用户
curl -X POST http://localhost:6116/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secure123"}'

# 登录获取 Token
curl -X POST http://localhost:6116/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secure123"}'

# 获取设备列表（使用 Token）
curl http://localhost:6116/api/devices \
  -H "Authorization: Bearer <your-token>"

# 添加设备
curl -X POST http://localhost:6116/api/devices \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"客厅智能灯","type":"smart_light"}'

# 发送控制指令
curl -X POST http://localhost:6116/api/control/send \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"device_id":"device-001","command":"ON"}'

# 健康检查
curl http://localhost:6116/health
```

---

## 容器管理

### 常用 Docker Compose 命令

```bash
# 查看运行状态
docker compose ps

# 查看实时日志
docker compose logs -f

# 重启服务
docker compose restart

# 停止服务
docker compose stop

# 停止并删除容器
docker compose down

# 更新到最新镜像
docker compose pull && docker compose up -d

# 查看镜像信息
docker images ghcr.io/mohan9898/iot
```

### 数据持久化

数据通过 Docker Volume `iot-manager-data` 持久化存储，映射到容器内 `/app/data` 目录：

- `iot.db` — 数据库文件
- `backup/` — 数据库备份目录

```bash
# 查看数据卷位置
docker volume inspect iot-manager-data

# 备份数据
docker cp iot-manager:/app/data ./data-backup

# 恢复数据
docker cp ./data-backup/. iot-manager:/app/data/
```

---

## MQTT 设备接入

### 支持的 Topic

| Topic | 方向 | 说明 |
|-------|------|------|
| `smart_light/#` | 订阅 | 智能灯设备所有消息 |
| `+/register` | 订阅 | 设备自动注册 |
| `+/status` | 订阅 | 设备状态上报 |
| `+/control` | 订阅 | 设备控制指令 |
| `+/metric` | 订阅 | 设备指标数据 |

### 设备自动注册

设备通过 MQTT 发送注册消息即可自动注册到平台：

```json
// Topic: <device_id>/register
{
  "id": "device-001",
  "name": "客厅智能灯",
  "type": "smart_light"
}
```

### 预设设备类型

系统预置了 7 种设备类型：

| 类型 ID | 名称 | 图标 | 支持指令 |
|---------|------|------|----------|
| `smart_light` | 智能灯 | 💡 | ON, OFF, AUTO, SET_THRESHOLD |
| `sensor` | 传感器 | 📡 | READ, CALIBRATE |
| `actuator` | 执行器 | ⚙️ | ON, OFF, TOGGLE |
| `controller` | 控制器 | 🔧 | START, STOP, RESET |
| `camera` | 摄像头 | 📷 | SNAPSHOT, RECORD, STREAM |
| `thermostat` | 恒温器 | 🌡️ | SET_TEMP, SET_MODE, OFF |
| `switch` | 智能开关 | 🔌 | ON, OFF, TOGGLE |

---

## 开发指南

### 环境要求

- Go 1.21+
- Node.js 18+（仅前端开发需要）
- MQTT Broker（推荐 EMQX，可选）

### 启动后端开发模式

```bash
cd backend

# 安装依赖
go mod download

# 设置环境变量
export JWT_SECRET="dev-secret-key"
export GIN_MODE="debug"
export SERVER_PORT=6116

# 启动服务
go run main.go
```

### 启动前端开发模式

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器（热更新）
npm run dev
# 默认运行在 http://localhost:5173
```

### 构建前端生产版本

```bash
cd frontend
npm run build
# 产物输出到 ../dist/ 目录
```

---

## 安全特性

- **非 root 运行** — Docker 容器以专用用户 `iot`（UID 1000）运行
- **JWT 认证** — bcrypt 密码哈希，支持 Token 过期管理
- **速率限制** — 每 IP 每分钟最多 100 次 API 请求
- **安全响应头** — 自动添加 CSP、X-Frame-Options、X-Content-Type-Options 等安全头
- **CORS 控制** — 默认同源策略，支持自定义允许源列表
- **无默认密钥** — JWT_SECRET 和 MQTT 凭证无硬编码默认值，必须通过环境变量配置
- **路径安全** — 数据库备份路径遍历保护

---

## 更新与维护

```bash
# 停止当前服务
docker compose down

# 拉取最新镜像
docker compose pull

# 重新启动
docker compose up -d

# 检查日志确认正常运行
docker compose logs -f
```

---

## 故障排除

### 服务无法启动

```bash
# 查看容器日志
docker compose logs iot-manager

# 检查端口占用
lsof -i :6116

# 检查 .env 文件是否正确
cat .env | grep JWT_SECRET
```

### MQTT 连接失败

系统支持 MQTT 优雅降级——即使 MQTT 代理不可用，Web 界面和设备管理功能仍可正常使用。

检查 MQTT 配置：
```bash
# 确认 MQTT 凭证已正确设置
echo $MQTT_USERNAME
echo $MQTT_PASSWORD

# 测试 MQTT 连接
mosquitto_sub -h <broker> -p <port> -u <username> -P <password> -t "#" -v
```

### 数据库问题

```bash
# 查看数据库文件
docker exec iot-manager ls -la /app/data/

# 查看备份
docker exec iot-manager ls -la /app/data/backup/

# 进入容器检查
docker exec -it iot-manager /bin/bash
```

---

## 许可证

MIT License