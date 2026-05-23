# IoT 设备管理系统

一个基于 Go + Vue 3 的智能物联网设备管理平台，支持设备注册、监控和远程控制。

## 📦 快速开始（推荐 - 使用 GitHub 镜像）

### 🚀 最简单的方式 - 一键部署脚本

```bash
# 1. 克隆仓库
git clone https://github.com/mohan9898/IOT.git
cd IOT

# 2. 运行快速部署脚本
chmod +x deploy.sh
./deploy.sh
```

### 📝 手动部署（核心版本）

```bash
# 1. 准备配置文件
cp .env.example .env
# 编辑 .env，配置您的 JWT_SECRET, MQTT_USERNAME, MQTT_PASSWORD

# 2. 拉取镜像并启动
docker-compose -f docker-compose-ghcr.yml up -d

# 3. 访问服务
# IoT Manager: http://localhost:6116
```

### 🌟 手动部署（完整版本 - 含监控）

```bash
# 1. 准备配置文件
cp .env.example .env
# 编辑 .env

# 2. 拉取镜像并启动完整服务
docker-compose -f docker-compose-ghcr-full.yml up -d

# 3. 访问服务
# IoT Manager: http://localhost:6116
# Grafana:     http://localhost:3000 (admin/admin)
# Prometheus:  http://localhost:9090
```

## 📚 文档导航

| 文档/脚本 | 说明 |
|-----------|------|
| [DEPLOYMENT.md](DEPLOYMENT.md) | 从源码构建部署教程 |
| [DOCKER_BUILD_GUIDE.md](DOCKER_BUILD_GUIDE.md) | GitHub Actions 自动构建使用说明 |
| [QUICKREFERENCE.md](QUICKREFERENCE.md) | 快速参考卡片 |
| [PRODUCTION_CHECKLIST.md](PRODUCTION_CHECKLIST.md) | 生产环境检查清单 |
| `./deploy.sh` | ⭐ **推荐使用** - 快速部署脚本（使用 GHCR 镜像） |
| `./check-status.sh` | 系统状态检查脚本 |
| `./cleanup.sh` | 停止和清理脚本 |

## 🏷️ GitHub Container Registry (GHCR) 镜像

项目使用 GitHub Actions 自动构建并发布 Docker 镜像：

```
ghcr.io/mohan9898/iot:latest
```

**可用标签**：
- `latest` - 最新版本
- `main` - main 分支构建
- `v1.0.0` - 版本号标签
- `sha-abc123` - 提交哈希标签

**拉取镜像**：
```bash
docker pull ghcr.io/mohan9898/iot:latest
```

## ✨ 功能特性

- 🔐 用户认证（JWT）
- 📱 响应式设计，支持移动端
- 💡 智能设备管理
- 📊 实时数据监控
- 🔌 设备远程控制
- 📡 MQTT 消息集成
- ⚡ WebSocket 实时通信
- 🛡️ 安全响应头
- 🔄 CORS 支持
- ⚖️ 速率限制
- 💾 自动数据备份
- 📜 结构化日志
- 📈 Prometheus 监控
- 📊 Grafana 仪表板
- 🚨 告警规则
- 🚀 GitHub Actions CI/CD 自动构建

## 🛠️ 技术栈

### 后端
- Go 1.21+
- Gin Web 框架
- SQLite 数据库
- Eclipse Paho MQTT
- JWT 认证
- Zap 日志
- Prometheus 指标

### 前端
- Vue 3
- Vite
- Tailwind CSS

### 监控
- Prometheus (指标采集)
- Grafana (可视化)
- Alertmanager (告警)

### CI/CD
- GitHub Actions
- GitHub Container Registry (GHCR)

## 🔧 配置说明

### 环境变量配置

主要配置通过 `.env` 文件：

```bash
# JWT 配置（重要！）
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRES_HOURS=24

# MQTT 配置
MQTT_BROKER=d11aab19.ala.cn-hangzhou.emqxsl.cn
MQTT_PORT=8883
MQTT_PROTOCOL=ssl
MQTT_TLS_ENABLED=true
MQTT_USERNAME=your-mqtt-user
MQTT_PASSWORD=your-mqtt-password

# 数据库配置
DB_BACKUP_ENABLE=true
DB_BACKUP_HOURS=24

# CORS 配置
CORS_ALLOWED_ORIGINS=*
```

### Docker Compose 配置文件

- `docker-compose-ghcr.yml` - 核心版本（仅 IoT Manager）
- `docker-compose-ghcr-full.yml` - 完整版本（含监控栈）
- `docker-compose.yml` - 从源码构建版本
- `docker-compose-core.yml` - 从源码构建核心版本

## 📊 监控与维护

### 更新到最新版本

```bash
# 拉取最新镜像
docker-compose -f docker-compose-ghcr.yml pull

# 重启服务
docker-compose -f docker-compose-ghcr.yml up -d
```

### 查看状态和日志

```bash
# 查看服务状态
docker-compose -f docker-compose-ghcr.yml ps

# 查看实时日志
docker-compose -f docker-compose-ghcr.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose-ghcr.yml logs -f iot-manager
```

### 数据备份和恢复

数据存储在 Docker volume `iot-manager-data` 中：

```bash
# 备份
docker run --rm -v iot-manager-data:/data -v $(pwd):/backup alpine tar czf /backup/iot-backup.tar.gz -C /data .

# 恢复
docker run --rm -v iot-manager-data:/data -v $(pwd):/backup alpine tar xzf /backup/iot-backup.tar.gz -C /data
```

## 📡 API 接口

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/update-account` - 更新账户
- `POST /api/auth/update-password` - 更新密码

### 设备接口
- `GET /api/devices` - 获取设备列表
- `POST /api/devices` - 创建设备
- `PUT /api/devices/:id` - 更新设备
- `DELETE /api/devices/:id` - 删除设备
- `GET /api/devices/stats` - 获取设备统计

### 控制接口
- `POST /api/control/send` - 发送控制命令
- `POST /api/control/threshold` - 设置阈值
- `GET /api/control/history/:id` - 获取历史记录

### 健康和监控
- `GET /health` - 服务健康状态
- `GET /metrics` - Prometheus 指标

## 🔒 安全特性

- ✅ JWT 令牌认证
- ✅ 密码 bcrypt 哈希
- ✅ MQTT over TLS
- ✅ 速率限制（每分钟 100 次）
- ✅ CORS 支持（可配置）
- ✅ 安全响应头
- ✅ 数据库自动备份
- ✅ GitHub Actions 安全构建流程

## 📈 开发说明

### 从源码构建和运行

```bash
# 后端开发
cd backend
go mod download
go run main.go

# 前端开发
cd frontend
npm install
npm run dev

# 从源码构建 Docker 镜像
docker-compose up -d --build
```

### 运行测试

```bash
cd backend
go test -v ./internal/db/...
```

## 🏛️ 项目结构

```
iot-manager/
├── backend/              # 后端代码
│   ├── config/          # 配置文件
│   ├── internal/        # 内部包
│   └── main.go          # 入口文件
├── frontend/            # 前端代码
├── monitoring/          # 监控配置
│   ├── prometheus/
│   └── grafana/
├── .github/             # GitHub 配置
│   └── workflows/       # GitHub Actions
├── .env.example         # 环境变量示例
├── deploy.sh            # ⭐ 快速部署脚本
├── docker-compose-ghcr.yml      # GHCR 核心版本
├── docker-compose-ghcr-full.yml # GHCR 完整版本
├── Dockerfile           # Docker 镜像构建
└── README.md
```

## 📋 阶段完成情况

- ✅ 阶段一（核心完善）- JWT 认证、HTTPS、数据备份、速率限制等
- ✅ 阶段二（监控运维）- Prometheus 指标、Grafana 仪表板、告警规则
- ✅ 阶段三（CI/CD）- GitHub Actions 自动构建、GHCR 镜像发布

## 🔖 部署方式对比

| 方式 | 说明 | 适用场景 |
|-----|------|---------|
| `./deploy.sh` | ⭐ 推荐 | 快速部署使用 |
| `docker-compose-ghcr.yml` | 核心版本 | 生产环境，仅需 IoT Manager |
| `docker-compose-ghcr-full.yml` | 完整版本 | 生产环境，需要监控 |
| `docker-compose.yml` | 源码构建 | 开发、自定义修改 |

## 📄 许可证

MIT
