# IoT Manager 核心系统部署指南

本文档说明如何只运行核心 IoT Manager 应用系统，不包含监控组件。

---

## 📦 系统组件说明

### IoT Manager（核心应用）
- **作用**：主要的物联网设备管理系统
- **功能**：用户认证、设备管理、设备控制、实时数据显示
- **技术栈**：Go + Vue 3 + SQLite + MQTT
- **端口**：6116

### Prometheus（可选监控）
- **作用**：指标采集和存储系统
- **功能**：采集 HTTP 请求数、设备状态、MQTT 连接等指标
- **端口**：9090
- **状态**：可选，仅用于监控

### Grafana（可选监控）
- **作用**：数据可视化仪表板
- **功能**：提供美观的监控图表和面板
- **端口**：3000
- **状态**：可选，仅用于监控

### Alertmanager（可选监控）
- **作用**：告警管理和通知
- **功能**：根据告警规则发送通知
- **端口**：9093
- **状态**：可选，仅用于监控

---

## 🚀 只运行核心系统

### 方法 1：使用核心配置文件（推荐）

```bash
# 使用专门的核心配置文件
docker-compose -f docker-compose-core.yml up -d

# 查看状态
docker-compose -f docker-compose-core.yml ps

# 查看日志
docker-compose -f docker-compose-core.yml logs -f

# 停止服务
docker-compose -f docker-compose-core.yml stop

# 停止并删除容器（保留数据）
docker-compose -f docker-compose-core.yml down
```

### 方法 2：使用快速启动脚本

```bash
# 运行快速启动脚本
./quickstart.sh

# 选择选项 2 - 仅部署 IoT Manager（轻量版）
```

### 方法 3：本地开发（无 Docker）

```bash
# 后端
cd backend
go mod download
go run main.go

# 前端（新终端）
cd frontend
npm install
npm run dev
```

---

## ✅ 核心系统功能完整度

即使不使用 Prometheus、Grafana、Alertmanager，IoT Manager 的以下功能仍然完整可用：

| 功能 | 可用性 | 说明 |
|------|---------|------|
| 用户认证（登录/注册） | ✅ 完整 | 与监控无关 |
| 设备管理 | ✅ 完整 | 与监控无关 |
| 设备控制 | ✅ 完整 | 与监控无关 |
| 实时数据显示 | ✅ 完整 | 通过 WebSocket，不依赖 Prometheus |
| MQTT 集成 | ✅ 完整 | 核心功能 |
| 数据备份 | ✅ 完整 | 内置功能 |
| JWT 认证 | ✅ 完整 | 核心功能 |
| 速率限制 | ✅ 完整 | 核心功能 |
| 安全响应头 | ✅ 完整 | 核心功能 |
| 监控仪表板 | ❌ 不可用 | 需要 Grafana |
| 告警通知 | ❌ 不可用 | 需要 Alertmanager |

---

## 📊 核心系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                      IoT Manager (Core)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │  Frontend    │  │   Backend    │  │     Database      │  │
│  │   (Vue 3)    │  │    (Go)      │  │    (SQLite)       │  │
│  └──────────────┘  └───────┬──────┘  └───────────────────┘  │
│                            │                                  │
│                      ┌──────▼──────┐                         │
│                      │  MQTT       │                         │
│                      │  (EMQX)     │                         │
│                      └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP (Port 6116)
                              │
                    ┌─────────▼─────────┐
                    │      User        │
                    └───────────────────┘
```

---

## 🔧 核心系统配置

### 环境变量

只需要配置核心系统需要的环境变量：

```env
# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=6116

# MQTT 配置（必须）
MQTT_BROKER=d11aab19.ala.cn-hangzhou.emqxsl.cn
MQTT_PORT=8883

# 数据库配置
DB_PATH=./data/iot.db
DB_BACKUP_ENABLE=true
DB_BACKUP_PATH=./data/backup

# 安全配置（重要！）
JWT_SECRET=your-strong-secret-key-here
```

### 无需配置的项

- Prometheus 相关配置：无需配置
- Grafana 相关配置：无需配置
- Alertmanager 相关配置：无需配置

---

## 📱 访问核心系统

部署成功后，只需访问一个地址：

- **IoT Manager**：http://localhost:6116

就这么简单！没有其他服务需要访问。

---

## 🧹 如果已运行完整版本，如何切换到核心版本

### 步骤 1：停止完整版本

```bash
# 停止完整版本（包含监控）
docker-compose down
```

### 步骤 2：启动核心版本

```bash
# 启动核心版本
docker-compose -f docker-compose-core.yml up -d
```

### 步骤 3：验证

```bash
# 查看容器状态（应该只有 iot-manager）
docker-compose -f docker-compose-core.yml ps

# 访问 http://localhost:6116 验证系统正常
```

---

## 💾 数据保留

两种部署方式使用相同的数据卷和数据结构，数据可以完美互用：

- ✅ 用户账户数据保留
- ✅ 设备信息数据保留
- ✅ 控制历史数据保留
- ✅ 从完整版本切换到核心版本：数据完整保留
- ✅ 从核心版本切换到完整版本：数据完整保留

---

## 🎯 什么时候使用核心版本

### 推荐使用核心版本的场景

- ✅ 生产环境初期部署
- ✅ 资源受限的环境（如小内存服务器）
- ✅ 不需要监控功能
- ✅ 简单的设备管理场景
- ✅ 快速测试和演示

### 推荐使用完整版本的场景

- ✅ 需要可视化监控
- ✅ 需要告警通知
- ✅ 大规模生产环境
- ✅ 需要性能和系统指标

---

## 🔄 在核心版本和完整版本之间切换

### 核心版本 → 完整版本

```bash
# 1. 停止核心版本
docker-compose -f docker-compose-core.yml down

# 2. 启动完整版本
docker-compose up -d
```

### 完整版本 → 核心版本

```bash
# 1. 停止完整版本
docker-compose down

# 2. 启动核心版本
docker-compose -f docker-compose-core.yml up -d
```

---

## 📝 总结

| 特性 | 核心版本 | 完整版本 |
|------|---------|---------|
| IoT Manager 功能 | ✅ 完整 | ✅ 完整 |
| Prometheus | ❌ | ✅ |
| Grafana | ❌ | ✅ |
| Alertmanager | ❌ | ✅ |
| 资源占用 | 低 | 中高 |
| 部署复杂度 | 简单 | 中等 |
| 监控能力 | 基础 | 强大 |

---

**选择适合你需求的部署方式！** 💪
