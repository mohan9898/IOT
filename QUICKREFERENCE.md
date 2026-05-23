# IoT Manager 快速参考卡

## 🔗 快速链接

| 内容 | 文件/链接 |
|------|----------|
| **核心系统部署**（推荐） | [CORE_DEPLOYMENT.md](CORE_DEPLOYMENT.md) |
| **完整部署教程** | [DEPLOYMENT.md](DEPLOYMENT.md) |
| **生产检查清单** | [PRODUCTION_CHECKLIST.md](PRODUCTION_CHECKLIST.md) |
| **主 README** | [README.md](README.md) |

---

## 🚀 快速命令

### 部署和管理

```bash
# 系统检查
./check-status.sh

# 快速启动（交互式）
./quickstart.sh

# 停止和清理
./cleanup.sh

# ========== 核心系统部署（推荐） ==========
# 只部署 IoT Manager（无监控）
docker-compose -f docker-compose-core.yml up -d

# 查看核心系统日志
docker-compose -f docker-compose-core.yml logs -f

# 停止核心系统
docker-compose -f docker-compose-core.yml stop

# ========== 完整系统部署 ==========
# 完整部署（含监控）
docker-compose up -d

# 查看完整系统日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f iot-manager
docker-compose logs -f prometheus
```

### 服务访问

#### 核心系统（推荐）

| 服务 | URL | 默认凭据 |
|------|-----|---------|
| IoT Manager | http://localhost:6116 | 自行注册 |

#### 完整系统（含监控）

| 服务 | URL | 默认凭据 |
|------|-----|---------|
| IoT Manager | http://localhost:6116 | 自行注册 |
| Prometheus | http://localhost:9090 | 无 |
| Grafana | http://localhost:3000 | admin/admin |
| Alertmanager | http://localhost:9093 | 无 |

---

## 📁 项目结构

```
iot-manager/
├── backend/              # Go 后端
├── frontend/             # Vue 3 前端
├── dist/                 # 构建产物
├── monitoring/           # 监控配置
│   ├── prometheus/       # Prometheus 配置
│   └── grafana/          # Grafana 配置
├── DEPLOYMENT.md         # 完整部署教程
├── PRODUCTION_CHECKLIST.md  # 生产检查清单
├── check-status.sh       # 系统状态检查
├── quickstart.sh         # 快速启动脚本
├── cleanup.sh            # 清理脚本
├── docker-compose.yml    # Docker Compose 配置
└── README.md             # 主文档
```

---

## ⚙️ 环境变量配置

复制 `.env.example` 为 `.env`，然后按需修改：

```env
# 服务器
SERVER_HOST=0.0.0.0
SERVER_PORT=6116

# HTTPS
HTTPS_ENABLE=false
HTTPS_CERT=./certs/cert.pem
HTTPS_KEY=./certs/key.pem

# MQTT
MQTT_BROKER=d11aab19.ala.cn-hangzhou.emqxsl.cn
MQTT_PORT=8883

# 数据库
DB_PATH=./data/iot.db
DB_BACKUP_ENABLE=false

# 安全
JWT_SECRET=your-very-secure-secret-key-change-in-production
```

---

## 🔒 生产环境关键修改

### 部署前必须修改

1. **JWT_SECRET** - 改为强随机密钥
2. **Grafana 密码** - 在 docker-compose.yml 中修改
3. **启用 HTTPS** - 使用 Let's Encrypt 等证书
4. **启用数据备份** - 设置 DB_BACKUP_ENABLE=true

---

## 🐛 常见问题速查

### 问题：容器无法启动

```bash
# 查看日志
docker-compose logs iot-manager

# 检查端口占用
netstat -tlnp | grep 6116
```

### 问题：MQTT 连接失败

1. 检查网络连接
2. 验证 EMQX 配置
3. 查看后端日志

### 问题：数据丢失

1. 检查数据卷是否正确挂载
2. 从备份恢复
3. 查看 DEPLOYMENT.md 中的备份章节

---

## 📊 Prometheus 指标

### HTTP 指标
- `http_requests_total` - 请求总数
- `http_request_duration_seconds` - 请求延迟

### 设备指标
- `devices_total` - 设备总数
- `devices_online` - 在线设备数
- `devices_offline` - 离线设备数

### MQTT 指标
- `mqtt_messages_total` - MQTT 消息数
- `mqtt_connection_status` - MQTT 连接状态

---

## 📞 获取帮助

1. 查看 DEPLOYMENT.md 中的详细教程
2. 查看 PRODUCTION_CHECKLIST.md 进行部署检查
3. 运行 `./check-status.sh` 查看系统状态

---

**祝使用愉快！** 🎉
