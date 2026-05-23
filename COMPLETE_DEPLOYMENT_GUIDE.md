# IoT 设备管理系统 - 完整详细部署指南

本文档提供了完整、详细、最新的部署流程，包含所有安全配置和优化。

---

## 目录

- [前置准备](#前置准备)
- [快速部署流程（推荐）](#快速部署流程推荐)
- [详细部署步骤](#详细部署步骤)
- [核心配置详解](#核心配置详解)
- [生产环境安全加固](#生产环境安全加固)
- [监控与维护](#监控与维护)
- [故障排查](#故障排查)

---

## 前置准备

### 系统要求

| 组件 | 最低版本 | 推荐版本 |
|------|---------|---------|
| Go | 1.20 | 1.21+ |
| Node.js | 16.x | 18.x+ |
| Docker | 20.10 | 24.0+ |
| Docker Compose | 2.0 | 2.20+ |

### 硬件要求

#### 开发环境
- CPU: 1核
- 内存: 1GB
- 存储: 10GB

#### 生产环境（最小）
- CPU: 2核
- 内存: 4GB
- 存储: 20GB SSD

#### 生产环境（推荐）
- CPU: 4核+
- 内存: 8GB+
- 存储: 50GB+ SSD

---

## 快速部署流程（推荐）

### 方式一：使用 Docker Compose（最快速）

```bash
# 1. 进入项目目录
cd /workspace/iot-manager

# 2. 复制并编辑环境变量文件
cp .env.example .env
# 重要！编辑 .env 文件，至少修改以下配置：
# - JWT_SECRET
# - MQTT_USERNAME
# - MQTT_PASSWORD

# 3. 启动服务（核心版，推荐）
docker-compose -f docker-compose-core.yml up -d

# 或者启动完整版本（含监控）
docker-compose up -d

# 4. 查看服务状态
docker-compose ps

# 5. 访问系统
# IoT Manager: http://localhost:6116
```

### 方式二：本地开发部署

```bash
# 后端
cd /workspace/iot-manager/backend
go mod download
go mod tidy
JWT_SECRET=your-strong-secret MQTT_USERNAME=your-mqtt-user MQTT_PASSWORD=your-mqtt-pass go run main.go

# 前端（新终端）
cd /workspace/iot-manager/frontend
npm install
npm run dev
```

---

## 详细部署步骤

### 第一步：环境准备

#### 1.1 检查系统环境

```bash
# 检查 Docker
docker --version
docker-compose --version

# 检查 Go（如果本地开发）
go version

# 检查 Node.js（如果本地开发）
node --version
npm --version
```

#### 1.2 获取项目代码

```bash
cd /workspace
# 如果使用 git
git clone <your-repo-url> iot-manager
cd iot-manager

# 或者使用现有项目
ls -la
```

#### 1.3 配置环境变量（重要！）

复制并编辑 `.env` 文件：

```bash
cp .env.example .env
```

**必须修改的配置项：**

| 配置项 | 说明 | 示例 |
|-------|------|------|
| `JWT_SECRET` | JWT 签名密钥 | `your-very-strong-secret-key-at-least-32-chars` |
| `MQTT_USERNAME` | MQTT 用户名 | `your-mqtt-username` |
| `MQTT_PASSWORD` | MQTT 密码 | `your-mqtt-password` |

**推荐配置项：**

| 配置项 | 说明 | 示例 |
|-------|------|------|
| `DB_BACKUP_ENABLE` | 启用自动备份 | `true` |
| `CORS_ALLOWED_ORIGINS` | 允许的 CORS 源 | `https://your-domain.com` |
| `SERVER_PORT` | 服务端口 | `6116` |

**完整的 .env 配置示例：**

```env
# ==========================================
# 服务器配置
# ==========================================
SERVER_HOST=0.0.0.0
SERVER_PORT=6116
HTTPS_ENABLE=false
HTTPS_CERT=
HTTPS_KEY=

# ==========================================
# MQTT 配置（重要！）
# ==========================================
MQTT_BROKER=d11aab19.ala.cn-hangzhou.emqxsl.cn
MQTT_PORT=8883
MQTT_PROTOCOL=ssl
MQTT_TLS_ENABLED=true
MQTT_CLIENT_ID_PREFIX=iot-manager-
MQTT_USERNAME=your-mqtt-username-here
MQTT_PASSWORD=your-mqtt-password-here

# ==========================================
# 数据库配置
# ==========================================
DB_PATH=./data/iot.db
DB_BACKUP_ENABLE=true
DB_BACKUP_PATH=./data/backup
DB_BACKUP_HOURS=24

# ==========================================
# JWT 安全配置（生产环境必须修改！）
# ==========================================
JWT_SECRET=your-very-secure-secret-key-change-this-in-production
JWT_EXPIRES_HOURS=24

# ==========================================
# CORS 配置（可选）
# ==========================================
CORS_ALLOWED_ORIGINS=*
# 或指定具体源：CORS_ALLOWED_ORIGINS=https://your-domain.com,http://localhost:5173

# ==========================================
# Grafana 配置（仅完整版本）
# ==========================================
GF_SECURITY_ADMIN_USER=admin
GF_SECURITY_ADMIN_PASSWORD=change-this-in-production
```

### 第二步：Docker 部署（推荐）

#### 2.1 选择部署版本

**核心版本（推荐用于生产，资源占用少）：**
- 只包含 IoT Manager 核心功能
- 适合生产环境部署
- 端口：仅 6116

**完整版本（含监控）：**
- 包含 IoT Manager + Prometheus + Grafana + Alertmanager
- 适合需要监控的环境
- 端口：6116, 3000, 9090, 9093

#### 2.2 部署核心版本

```bash
cd /workspace/iot-manager

# 1. 创建数据目录
mkdir -p data/backup
chmod 755 data

# 2. 构建并启动
docker-compose -f docker-compose-core.yml up -d --build

# 3. 查看启动日志
docker-compose -f docker-compose-core.yml logs -f

# 4. 等待服务启动，检查健康状态
curl http://localhost:6116/health
```

#### 2.3 部署完整版本（含监控）

```bash
cd /workspace/iot-manager

# 1. 检查监控配置
ls -la monitoring/
ls -la monitoring/prometheus
ls -la monitoring/grafana

# 2. 构建并启动所有服务
docker-compose up -d --build

# 3. 查看所有服务状态
docker-compose ps

# 4. 查看实时日志
docker-compose logs -f

# 5. 检查各个服务
curl http://localhost:6116/health
curl http://localhost:9090/-/healthy
curl http://localhost:3000/api/health
```

### 第三步：验证部署

#### 3.1 检查服务状态

```bash
# 查看所有容器
docker-compose ps

# 预期输出（核心版本）：
# NAME                STATUS
# iot-manager         Up 2 minutes (healthy)

# 预期输出（完整版本）：
# NAME                STATUS
# iot-manager         Up 2 minutes (healthy)
# prometheus          Up 2 minutes
# grafana             Up 2 minutes
# alertmanager        Up 2 minutes
```

#### 3.2 测试 API 接口

```bash
# 健康检查
curl -s http://localhost:6116/health | python3 -m json.tool

# 指标端点
curl -s http://localhost:6116/metrics | head -20

# 测试用户注册
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}' \
  http://localhost:6116/api/auth/register

# 测试登录
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}' \
  http://localhost:6116/api/auth/login
```

#### 3.3 访问 Web 界面

打开浏览器访问：

| 服务 | URL | 默认凭据 |
|------|-----|---------|
| IoT Manager | http://localhost:6116 | 自行注册 |
| Prometheus | http://localhost:9090 | 无 |
| Grafana | http://localhost:3000 | admin/admin |
| Alertmanager | http://localhost:9093 | 无 |

---

## 核心配置详解

### 后端配置文件

主要配置位于 [`backend/config/config.go`](file:///workspace/iot-manager/backend/config/config.go)，但推荐通过环境变量配置。

### 环境变量完整列表

#### 服务器配置

| 变量名 | 类型 | 默认值 | 说明 |
|-------|------|-------|------|
| `SERVER_HOST` | string | `0.0.0.0` | 监听地址 |
| `SERVER_PORT` | int | `6116` | 监听端口 |
| `HTTPS_ENABLE` | bool | `false` | 启用 HTTPS |
| `HTTPS_CERT` | string | `` | HTTPS 证书路径 |
| `HTTPS_KEY` | string | `` | HTTPS 私钥路径 |
| `CORS_ALLOWED_ORIGINS` | string | `` | CORS 允许的源，逗号分隔 |

#### MQTT 配置

| 变量名 | 类型 | 默认值 | 说明 |
|-------|------|-------|------|
| `MQTT_BROKER` | string | `d11aab19.ala.cn-hangzhou.emqxsl.cn` | MQTT Broker 地址 |
| `MQTT_PORT` | int | `8883` | MQTT 端口 |
| `MQTT_PROTOCOL` | string | `ssl` | 协议（`ssl`/`tcp`/`ws`/`wss`） |
| `MQTT_TLS_ENABLED` | bool | `true` | 启用 TLS |
| `MQTT_USERNAME` | string | `` | MQTT 用户名 |
| `MQTT_PASSWORD` | string | `` | MQTT 密码 |
| `MQTT_CLIENT_ID_PREFIX` | string | `iot-manager-` | 客户端 ID 前缀 |

#### 数据库配置

| 变量名 | 类型 | 默认值 | 说明 |
|-------|------|-------|------|
| `DB_PATH` | string | `./data/iot.db` | 数据库文件路径 |
| `DB_BACKUP_ENABLE` | bool | `false` | 启用自动备份 |
| `DB_BACKUP_PATH` | string | `./data/backup` | 备份目录 |
| `DB_BACKUP_HOURS` | int | `24` | 备份间隔（小时） |

#### JWT 安全配置

| 变量名 | 类型 | 默认值 | 说明 |
|-------|------|-------|------|
| `JWT_SECRET` | string | `` | JWT 签名密钥（重要！） |
| `JWT_EXPIRES_HOURS` | int | `24` | Token 有效期（小时） |

---

## 生产环境安全加固

### 1. 修改所有默认密码

#### JWT_SECRET 配置

```bash
# 生成强密钥
openssl rand -base64 32

# 编辑 .env 文件
JWT_SECRET=your-generated-strong-secret-key
```

#### Grafana 密码

编辑 [`docker-compose.yml`](file:///workspace/iot-manager/docker-compose.yml)：

```yaml
grafana:
  environment:
    - GF_SECURITY_ADMIN_USER=your-admin-username
    - GF_SECURITY_ADMIN_PASSWORD=your-strong-admin-password
```

### 2. 配置 HTTPS（强烈推荐）

#### 使用 Let's Encrypt 免费证书

```bash
# 1. 安装 Certbot
apt update && apt install -y certbot

# 2. 获取证书（需要域名解析）
certbot certonly --standalone -d your-domain.com

# 3. 证书位置
# /etc/letsencrypt/live/your-domain.com/fullchain.pem
# /etc/letsencrypt/live/your-domain.com/privkey.pem
```

#### 更新 Docker Compose 配置

修改 [`docker-compose-core.yml`](file:///workspace/iot-manager/docker-compose-core.yml)：

```yaml
services:
  iot-manager:
    environment:
      - HTTPS_ENABLE=true
      - HTTPS_CERT=/etc/letsencrypt/live/your-domain.com/fullchain.pem
      - HTTPS_KEY=/etc/letsencrypt/live/your-domain.com/privkey.pem
    volumes:
      - /etc/letsencrypt:/etc/letsencrypt:ro
```

#### 配置自动续期

```bash
# 添加定时任务
crontab -e

# 添加以下行（每月 1 号凌晨 2 点续期）
0 2 1 * * certbot renew --quiet && docker-compose -f docker-compose-core.yml restart
```

### 3. 配置防火墙

#### UFW（Ubuntu/Debian）

```bash
# 允许 SSH
ufw allow 22/tcp

# 允许 IoT Manager
ufw allow 6116/tcp

# 如果使用 HTTPS
ufw allow 443/tcp
ufw allow 80/tcp

# 如果需要从外网访问监控（谨慎）
# ufw allow 3000/tcp  # Grafana
# ufw allow 9090/tcp  # Prometheus
# ufw allow 9093/tcp  # Alertmanager

# 启用防火墙
ufw enable

# 查看状态
ufw status
```

#### Firewalld（CentOS/RHEL）

```bash
firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-port=6116/tcp
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload
```

### 4. 配置 Nginx 反向代理（可选但推荐）

```nginx
# /etc/nginx/sites-available/iot-manager
server {
    listen 80;
    server_name your-domain.com;
    
    # 强制 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    # SSL 配置
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # IoT Manager 代理
    location / {
        proxy_pass http://localhost:6116;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
    
    # WebSocket 支持
    location /api/ws {
        proxy_pass http://localhost:6116/api/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # Grafana 代理（如果启用）
    location /grafana/ {
        proxy_pass http://localhost:3000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

启用 Nginx 配置：

```bash
ln -s /etc/nginx/sites-available/iot-manager /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx
```

### 5. 限制监控服务访问

编辑 [`docker-compose.yml`](file:///workspace/iot-manager/docker-compose.yml)，限制监控服务只监听 localhost：

```yaml
services:
  prometheus:
    ports:
      - "127.0.0.1:9090:9090"
  
  grafana:
    ports:
      - "127.0.0.1:3000:3000"
  
  alertmanager:
    ports:
      - "127.0.0.1:9093:9093"
```

### 6. 配置 Docker 资源限制

编辑 [`docker-compose.yml`](file:///workspace/iot-manager/docker-compose.yml)：

```yaml
services:
  iot-manager:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
  
  prometheus:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
  
  grafana:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
```

---

## 监控与维护

### 日常检查清单

- [ ] 检查容器状态：`docker-compose ps`
- [ ] 查看系统日志：`docker-compose logs --tail=100`
- [ ] 检查磁盘空间：`df -h`
- [ ] 检查内存使用：`free -h`
- [ ] 检查备份文件：`ls -la data/backup/`
- [ ] 访问 Grafana 监控面板

### 定期维护任务

| 任务 | 频率 | 操作 |
|------|------|------|
| 数据备份验证 | 每周 | 测试恢复备份数据 |
| 系统更新 | 每月 | 更新 Docker 镜像和系统包 |
| 安全扫描 | 每月 | 运行安全扫描工具 |
| 日志清理 | 每月 | 清理旧日志文件 |
| 性能检查 | 每季度 | 分析系统性能瓶颈 |

### 常用运维命令

```bash
# 查看服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f iot-manager
docker-compose logs --tail=200 iot-manager

# 启动/停止/重启服务
docker-compose start
docker-compose stop
docker-compose restart

# 更新服务
docker-compose pull
docker-compose up -d

# 重新构建并部署
docker-compose up -d --build

# 进入容器
docker-compose exec iot-manager sh

# 清理未使用的资源
docker system prune -a
```

### 数据备份与恢复

#### 手动触发备份

```bash
# 如果备份管理器正在运行，会自动定期备份
# 手动备份可以直接复制数据库文件
cp data/iot.db data/backup/iot-manual-$(date +%Y%m%d-%H%M%S).db
```

#### 从备份恢复

```bash
# 1. 停止服务
docker-compose stop

# 2. 备份当前数据库（以防万一）
cp data/iot.db data/iot.db.backup

# 3. 恢复备份
cp data/backup/your-backup-file.db data/iot.db

# 4. 设置正确权限
chmod 644 data/iot.db

# 5. 启动服务
docker-compose start
```

---

## 故障排查

### 问题 1：容器无法启动

```bash
# 查看容器日志
docker-compose logs iot-manager

# 常见原因：
# - 端口被占用
# - 权限问题
# - 配置错误

# 检查端口占用
netstat -tlnp | grep 6116

# 或者使用 lsof
lsof -i :6116
```

### 问题 2：MQTT 连接失败

应用程序现在可以在 MQTT 连接失败时继续运行，但控制功能会受影响。

```bash
# 检查 MQTT 配置
docker-compose exec iot-manager env | grep MQTT

# 测试网络连接
telnet d11aab19.ala.cn-hangzhou.emqxsl.cn 8883

# 检查防火墙
ufw status
```

### 问题 3：健康检查失败

```bash
# 检查健康状态
curl -v http://localhost:6116/health

# 查看应用日志
docker-compose logs iot-manager --tail=50

# 检查数据库文件
ls -la data/
```

### 问题 4：Grafana 无法访问

```bash
# 检查 Grafana 容器状态
docker-compose ps grafana

# 查看 Grafana 日志
docker-compose logs grafana

# 重置 Grafana（谨慎使用，会丢失自定义配置）
docker-compose down -v
docker-compose up -d
```

### 问题 5：数据无法写入

```bash
# 检查数据目录权限
ls -la data/

# 修复权限
chmod 755 data
chmod 644 data/iot.db

# 检查磁盘空间
df -h
```

### 获取帮助

如遇到问题，请按以下顺序排查：

1. 查看容器日志：`docker-compose logs`
2. 检查系统资源：`free -h`, `df -h`
3. 检查网络连接：`ping`, `telnet`
4. 参考常见问题章节
5. 联系技术支持

---

## 性能优化建议

### 数据库优化

- 使用 SSD 存储数据库文件
- 定期执行 SQLite VACUUM
- 配置合理的缓存大小

### Docker 优化

- 使用官方镜像
- 配置资源限制
- 启用 Docker 日志轮转

### 系统优化

- 启用 swap（如果内存紧张）
- 配置文件描述符限制
- 优化网络参数

---

## 升级流程

### 从旧版本升级

```bash
# 1. 备份数据
cp -r data data.backup.$(date +%Y%m%d)

# 2. 拉取最新代码
git pull

# 3. 查看变更
git log --oneline -10

# 4. 停止旧版本
docker-compose down

# 5. 重新构建并启动
docker-compose up -d --build

# 6. 验证服务正常
docker-compose ps
curl http://localhost:6116/health
```

---

## 附录

### A. 完整端口列表

| 端口 | 服务 | 协议 | 说明 |
|------|------|------|------|
| 6116 | IoT Manager | TCP | 主服务端口 |
| 3000 | Grafana | TCP | 监控仪表板 |
| 9090 | Prometheus | TCP | 指标采集 |
| 9093 | Alertmanager | TCP | 告警管理 |
| 8883 | EMQX MQTT | TCP/S | MQTT Broker（外部） |

### B. 目录结构

```
iot-manager/
├── backend/                 # 后端代码
│   ├── config/             # 配置文件
│   ├── internal/           # 内部包
│   ├── data/               # 数据库（运行时）
│   └── main.go             # 入口文件
├── frontend/               # 前端代码
├── monitoring/             # 监控配置
│   ├── prometheus/         # Prometheus 配置
│   └── grafana/            # Grafana 配置
├── data/                   # 数据目录（运行时）
│   ├── iot.db              # 主数据库
│   └── backup/             # 备份文件
├── .env                    # 环境变量
├── .env.example            # 环境变量示例
├── docker-compose.yml      # 完整版本配置
├── docker-compose-core.yml # 核心版本配置
└── Dockerfile              # Docker 镜像构建
```

---

## 总结

本部署指南涵盖了从快速开始到生产环境配置的完整流程。按照以下步骤操作：

1. **准备环境**：检查系统要求，获取代码
2. **配置环境变量**：重点配置 JWT_SECRET、MQTT 凭据
3. **选择部署方式**：核心版本或完整版本
4. **启动服务**：使用 Docker Compose
5. **验证部署**：检查服务状态，测试 API
6. **安全加固**：配置 HTTPS、防火墙、密码
7. **监控维护**：建立日常检查和维护流程

祝部署顺利！如有问题请参考故障排查章节。

---

**最后更新时间**：2026-05-23
**文档版本**：2.0.0
