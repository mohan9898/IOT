# IoT 设备管理系统

一个基于 Go + Vue 3 的智能物联网设备管理平台，支持设备注册、监控和远程控制。

## 📚 快速开始

| 文档/脚本 | 说明 |
|-----------|------|
| [CORE_DEPLOYMENT.md](CORE_DEPLOYMENT.md) | 核心系统部署（**推荐**，仅 IoT Manager） |
| [QUICKREFERENCE.md](QUICKREFERENCE.md) | 快速参考卡片 |
| [DEPLOYMENT.md](DEPLOYMENT.md) | 完整部署教程（含监控） |
| [PRODUCTION_CHECKLIST.md](PRODUCTION_CHECKLIST.md) | 生产环境检查清单 |
| `./check-status.sh` | 系统状态检查脚本 |
| `./quickstart.sh` | 交互式快速部署脚本 |
| `./cleanup.sh` | 停止和清理脚本 |

### 最快速开始（仅核心系统，推荐）

```bash
# 1. 检查系统状态
./check-status.sh

# 2. 启动核心系统
docker-compose -f docker-compose-core.yml up -d

# 3. 访问服务
# IoT Manager: http://localhost:6116
```

### 完整部署（含监控）

```bash
# 1. 检查系统状态
./check-status.sh

# 2. 快速部署
./quickstart.sh

# 3. 访问服务
# IoT Manager: http://localhost:6116
# Grafana: http://localhost:3000 (admin/admin)
```

## 功能特性

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

## 技术栈

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

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- MQTT Broker（默认使用 EMQX Cloud）
- Docker 和 Docker Compose (可选，用于监控栈)

### 本地开发

#### 环境配置
复制环境变量示例文件：
```bash
cp .env.example .env
# 根据需要编辑 .env 文件
```

#### 后端
```bash
cd backend
go mod download
go run main.go
```

#### 前端
```bash
cd frontend
npm install
npm run dev
```

#### 运行测试
```bash
cd backend
go test -v ./internal/db/...
```

### Docker 部署（含监控）

#### 完整部署（含监控）
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

#### 访问地址
- IoT Manager: http://localhost:6116
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (用户名/密码: admin/admin)
- Alertmanager: http://localhost:9093

## 配置说明

### 服务器配置
编辑 `backend/config/config.go` 或使用环境变量：
- `SERVER_HOST`: 监听地址，默认 0.0.0.0
- `SERVER_PORT`: 监听端口，默认 6116
- `HTTPS_ENABLE`: 启用 HTTPS，默认 false
- `HTTPS_CERT`: 证书路径
- `HTTPS_KEY`: 私钥路径

### MQTT 配置
- `MQTT_BROKER`: MQTT 代理地址
- `MQTT_PORT`: MQTT 端口
- 用户名和密码在代码中配置

### 数据库配置
- `DB_PATH`: 数据库文件路径
- `DB_BACKUP_ENABLE`: 启用自动备份，默认 false
- `DB_BACKUP_PATH`: 备份路径

### JWT 配置
- `JWT_SECRET`: JWT 签名密钥（生产环境必须修改！）

## API 文档

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

### 健康检查
- `GET /health` - 服务健康状态

### 监控接口
- `GET /metrics` - Prometheus 指标

## Prometheus 指标

系统暴露以下 Prometheus 指标：

### HTTP 指标
- `http_requests_total`: HTTP 请求总数
- `http_request_duration_seconds`: HTTP 请求延迟分布

### 设备指标
- `devices_total`: 设备总数
- `devices_online`: 在线设备数
- `devices_offline`: 离线设备数

### MQTT 指标
- `mqtt_messages_total`: MQTT 消息总数
- `mqtt_connection_status`: MQTT 连接状态

### 控制指标
- `control_commands_total`: 控制命令总数

### 认证指标
- `login_attempts_total`: 登录尝试总数
- `active_users`: 活跃用户数

## 安全特性

- ✅ JWT 令牌认证
- ✅ 密码 bcrypt 哈希
- ✅ MQTT over TLS
- ✅ 速率限制（每分钟 100 次）
- ✅ CORS 支持
- ✅ 安全响应头（X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, CSP）
- ✅ 数据库自动备份（可配置）

## 告警规则

系统配置了以下告警规则：

- `ServiceDown`: 服务不可用（严重）
- `HighOfflineDevices`: 超过 30% 设备离线（警告）
- `MQTTDisconnected`: MQTT 连接断开（严重）
- `HighHTTPErrorRate`: HTTP 错误率 > 10%（警告）
- `HighLoginFailureRate`: 登录失败率过高（警告）
- `HighRequestLatency`: 请求延迟过高（警告）

## 项目结构

```
iot-manager/
├── backend/
│   ├── config/
│   ├── internal/
│   │   ├── api/
│   │   ├── auth/
│   │   ├── db/
│   │   ├── mqtt/
│   │   └── metrics/
│   ├── go.mod
│   └── main.go
├── frontend/
├── monitoring/
│   ├── prometheus/
│   │   ├── prometheus.yml
│   │   ├── alertmanager.yml
│   │   └── rules/
│   └── grafana/
│       └── provisioning/
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 阶段一（核心完善）- ✅ 已完成

- ✅ JWT 认证
- ✅ 数据库操作测试
- ✅ HTTPS 支持
- ✅ 数据备份机制
- ✅ 结构化日志
- ✅ 速率限制
- ✅ 安全响应头
- ✅ CORS 支持
- ✅ Docker 部署

## 阶段二（监控运维）- ✅ 已完成

- ✅ Prometheus 指标暴露
- ✅ Grafana 仪表板
- ✅ 告警规则配置
- ✅ 监控栈 Docker Compose

## 生产部署建议

1. **HTTPS**: 配置 SSL/TLS 证书（Let's Encrypt 推荐）
2. **数据备份**: 启用自动备份，定期验证备份文件
3. **监控告警**: 使用提供的 Prometheus + Grafana 监控栈
4. **日志收集**: 可以扩展使用 ELK 或类似方案
5. **防火墙**: 限制访问端口
6. **强密码**: 修改 JWT_SECRET、Grafana 密码等为强随机密钥

## 许可证

MIT
