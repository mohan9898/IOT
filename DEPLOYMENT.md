# IoT 设备管理系统 - 完整部署教程

本文档提供了多种部署方式的详细指南，包括本地开发、Docker 部署、云平台部署等。

## 目录

- [前置准备](#前置准备)
- [方式一：本地开发部署](#方式一本地开发部署)
- [方式二：Docker 单机部署](#方式二docker-单机部署)
- [方式三：Docker Compose 完整部署（推荐）](#方式三docker-compose-完整部署推荐)
- [方式四：云平台部署](#方式四云平台部署)
- [部署后配置](#部署后配置)
- [常见问题](#常见问题)

---

## 前置准备

### 系统要求

| 组件 | 最低版本 | 推荐版本 |
|------|---------|---------|
| Go | 1.20 | 1.21+ |
| Node.js | 16.x | 18.x+ |
| Docker | 20.10 | 24.0+ |
| Docker Compose | 2.0 | 2.20+ |

### 准备工作

1. **获取项目代码**

```bash
# 如果是 git 仓库
git clone <your-repo-url>
cd iot-manager

# 或者直接使用现有项目
cd /workspace/iot-manager
```

2. **配置环境变量**

```bash
# 复制示例配置文件
cp .env.example .env

# 根据需要编辑 .env 文件
# 修改 JWT_SECRET、Grafana 密码等敏感信息
```

---

## 方式一：本地开发部署

适合开发调试，不适合生产环境。

### 步骤 1：后端部署

```bash
cd backend

# 安装依赖
go mod download
go mod tidy

# 验证 MQTT 配置
# 编辑 config/config.go 确认 MQTT 配置正确

# 运行后端
go run main.go
```

后端将在 http://localhost:6116 启动。

### 步骤 2：前端部署

```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行（热重载）
npm run dev
```

前端将在 http://localhost:5173 启动。

### 步骤 3：构建生产版本前端

```bash
cd frontend

# 构建生产版本
npm run build

# 构建产物会输出到 dist 目录
# 后端会自动服务 dist 目录的内容
```

### 步骤 4：验证部署

1. 访问 http://localhost:6116 确认系统正常
2. 创建测试账户
3. 登录并测试设备管理功能

---

## 方式二：Docker 单机部署

只部署 IoT Manager 核心服务，不包含监控栈。

### 步骤 1：构建 Docker 镜像

```bash
cd /workspace/iot-manager

# 构建镜像
docker build -t iot-manager:latest .

# 验证镜像
docker images | grep iot-manager
```

### 步骤 2：创建数据目录

```bash
# 创建数据持久化目录
mkdir -p ./data
chmod 755 ./data
```

### 步骤 3：运行容器

```bash
# 方式 A：使用命令行直接运行
docker run -d \
  --name iot-manager \
  -p 6116:6116 \
  -v $(pwd)/data:/root/data \
  -e MQTT_BROKER=d11aab19.ala.cn-hangzhou.emqxsl.cn \
  -e MQTT_PORT=8883 \
  --restart unless-stopped \
  iot-manager:latest

# 方式 B：使用 docker-compose（轻量版）
cat > docker-compose-simple.yml << 'EOF'
version: '3.8'
services:
  iot-manager:
    image: iot-manager:latest
    container_name: iot-manager
    ports:
      - "6116:6116"
    volumes:
      - ./data:/root/data
    environment:
      - MQTT_BROKER=d11aab19.ala.cn-hangzhou.emqxsl.cn
      - MQTT_PORT=8883
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:6116/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
EOF

# 启动轻量版
docker-compose -f docker-compose-simple.yml up -d
```

### 步骤 4：验证部署

```bash
# 查看容器状态
docker ps

# 查看日志
docker logs -f iot-manager

# 测试健康检查
curl http://localhost:6116/health
```

---

## 方式三：Docker Compose 完整部署（推荐）

包含 IoT Manager、Prometheus、Grafana、Alertmanager 的完整监控栈。

### 步骤 1：准备环境

```bash
cd /workspace/iot-manager

# 确保配置文件存在
ls -la docker-compose.yml
ls -la monitoring/
```

### 步骤 2：配置环境变量

编辑 `.env` 文件：

```bash
# 复制并编辑
cp .env.example .env

# 建议修改的配置
JWT_SECRET=your-very-secure-secret-key-change-this-in-production
DB_BACKUP_ENABLE=true
```

### 步骤 3：启动完整服务

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 或者只启动不构建
docker-compose up -d
```

### 步骤 4：验证服务启动

```bash
# 查看所有容器状态
docker-compose ps

# 应该看到 4 个容器：
# - iot-manager
# - prometheus
# - grafana
# - alertmanager

# 查看日志（实时）
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f iot-manager
docker-compose logs -f prometheus
```

### 步骤 5：访问各个服务

打开浏览器访问以下地址：

| 服务 | URL | 默认账号 |
|------|-----|---------|
| IoT Manager | http://localhost:6116 | 自行注册 |
| Prometheus | http://localhost:9090 | 无 |
| Grafana | http://localhost:3000 | admin/admin |
| Alertmanager | http://localhost:9093 | 无 |

### 步骤 6：配置 Grafana（自动完成）

Grafana 已自动配置：
- ✅ 数据源自动连接到 Prometheus
- ✅ 仪表板自动导入
- ✅ 可以直接查看监控面板

### 步骤 7：常用命令

```bash
# 停止服务
docker-compose stop

# 启动服务
docker-compose start

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 更新服务（重新构建）
docker-compose up -d --build

# 删除所有服务和数据
docker-compose down -v
```

---

## 方式四：云平台部署

### 部署到阿里云 ECS

#### 准备

1. 购买阿里云 ECS 实例（推荐 2核4G 以上）
2. 配置安全组开放端口：6116, 3000, 9090, 9093
3. 获取服务器公网 IP

#### 部署步骤

```bash
# SSH 连接服务器
ssh root@your-ecs-ip

# 安装 Docker 和 Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 启动 Docker
systemctl start docker
systemctl enable docker

# 安装 Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 上传项目代码到服务器
# 可以使用 scp 或 git clone

# 进入项目目录
cd /root/iot-manager

# 配置环境变量
cp .env.example .env
# 编辑 .env，修改 JWT_SECRET 等

# 启动服务
docker-compose up -d

# 配置反向代理（可选，使用 Nginx）
apt update && apt install -y nginx

cat > /etc/nginx/sites-available/iot-manager << 'EOF'
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:6116;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    location /grafana/ {
        proxy_pass http://localhost:3000/;
        proxy_set_header Host $host;
    }
}
EOF

ln -s /etc/nginx/sites-available/iot-manager /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx

# 配置防火墙
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

### 部署到腾讯云 CVM

步骤类似阿里云，请参考阿里云部署文档。

### 部署到 AWS EC2

```bash
# SSH 连接
ssh -i your-key.pem ec2-user@your-ec2-ip

# 安装 Docker
sudo yum update -y
sudo amazon-linux-extras install docker
sudo service docker start
sudo usermod -a -G docker ec2-user

# 后续步骤同阿里云
```

### 部署到 Kubernetes（高级）

创建 Kubernetes 部署文件：

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: iot-manager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: iot-manager
  template:
    metadata:
      labels:
        app: iot-manager
    spec:
      containers:
      - name: iot-manager
        image: iot-manager:latest
        ports:
        - containerPort: 6116
        env:
        - name: MQTT_BROKER
          value: "d11aab19.ala.cn-hangzhou.emqxsl.cn"
        - name: MQTT_PORT
          value: "8883"
        volumeMounts:
        - name: data
          mountPath: /root/data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: iot-data
---
apiVersion: v1
kind: Service
metadata:
  name: iot-manager
spec:
  type: LoadBalancer
  selector:
    app: iot-manager
  ports:
  - protocol: TCP
    port: 6116
    targetPort: 6116
```

部署命令：

```bash
kubectl apply -f k8s-deployment.yaml
```

---

## 部署后配置

### 1. 修改默认密码

#### 修改 JWT_SECRET

编辑 [`backend/config/config.go`](file:///workspace/iot-manager/backend/config/config.go)：

```go
JWT: JWTConfig{
    Secret:       "your-new-very-secure-secret-key-here-make-it-long",
    ExpiresHours: 24,
},
```

#### 修改 Grafana 密码

编辑 [`docker-compose.yml`](file:///workspace/iot-manager/docker-compose.yml)：

```yaml
grafana:
  environment:
    - GF_SECURITY_ADMIN_USER=your-username
    - GF_SECURITY_ADMIN_PASSWORD=your-strong-password
```

重启服务：

```bash
docker-compose up -d
```

### 2. 配置 HTTPS（推荐）

#### 使用 Let's Encrypt（免费）

```bash
# 安装 Certbot
apt update
apt install -y certbot python3-certbot-nginx

# 获取证书
certbot --nginx -d your-domain.com

# 自动续期
certbot renew --dry-run
```

#### 更新 Nginx 配置

```nginx
server {
    listen 443 ssl;
    listen [::]:443 ssl;
    
    server_name your-domain.com;
    
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    location / {
        proxy_pass http://localhost:6116;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

### 3. 启用数据备份

编辑 `.env` 文件：

```env
DB_BACKUP_ENABLE=true
DB_BACKUP_PATH=./data/backup
```

### 4. 配置告警通知（可选）

编辑 [`monitoring/prometheus/alertmanager.yml`](file:///workspace/iot-manager/monitoring/prometheus/alertmanager.yml)：

```yaml
receivers:
  - name: 'default-receiver'
    webhook_configs:
    - url: 'https://your-webhook-url/alert'
    email_configs:
    - to: 'your-email@example.com'
      from: 'alert@your-domain.com'
      smarthost: 'smtp.example.com:587'
      auth_username: 'your-smtp-user'
      auth_password: 'your-smtp-password'
```

### 5. 配置防火墙

```bash
# UFW (Ubuntu/Debian)
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw allow 6116/tcp  # IoT Manager
ufw enable

# Firewalld (CentOS/RHEL)
firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-port=6116/tcp
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload
```

---

## 常见问题

### Q1: Docker 容器无法启动？

```bash
# 检查日志
docker-compose logs iot-manager

# 常见问题：端口被占用
netstat -tlnp | grep 6116
# 修改 docker-compose.yml 中的端口映射
```

### Q2: MQTT 连接失败？

1. 检查网络连接
2. 验证 MQTT 配置（地址、端口、用户名、密码）
3. 检查防火墙设置
4. 查看容器日志

### Q3: 数据丢失？

1. 确认数据卷正确挂载
2. 检查备份是否启用
3. 从备份恢复（如需要）

```bash
# 查看备份文件
ls -la ./data/backup/

# 恢复备份（参考 backup.go 中的 Restore 方法）
```

### Q4: Grafana 无法访问？

```bash
# 检查 Grafana 容器状态
docker-compose ps grafana

# 查看 Grafana 日志
docker-compose logs grafana

# 重置 Grafana 数据（谨慎使用）
docker-compose down -v
docker-compose up -d
```

### Q5: 如何升级版本？

```bash
# 拉取最新代码
git pull

# 重新构建并部署
docker-compose down
docker-compose up -d --build

# 检查服务状态
docker-compose ps
```

### Q6: 性能优化建议？

1. 使用 SSD 存储数据库
2. 配置 Docker 资源限制
3. 启用 HTTPS 缓存
4. 配置 Prometheus 数据保留策略
5. 定期清理旧备份文件

---

## 维护操作手册

### 日常检查清单

- [ ] 检查容器运行状态 `docker-compose ps`
- [ ] 查看系统日志 `docker-compose logs --tail=100`
- [ ] 检查 Grafana 监控面板
- [ ] 验证数据备份正常
- [ ] 检查磁盘空间 `df -h`

### 定期维护任务

| 任务 | 频率 | 操作 |
|------|------|------|
| 数据备份验证 | 每周 | 测试恢复备份数据 |
| 系统更新 | 每月 | 更新 Docker 镜像和系统包 |
| 安全扫描 | 每月 | 运行安全扫描工具 |
| 日志清理 | 每月 | 清理旧日志文件 |
| 性能检查 | 每季度 | 分析系统性能瓶颈 |

### 紧急恢复流程

1. 停止服务：`docker-compose stop`
2. 从备份恢复数据
3. 检查问题原因
4. 修复问题
5. 重新启动服务：`docker-compose start`
6. 验证服务正常

---

## 联系支持

如遇到问题，请：

1. 查看日志文件
2. 检查监控面板
3. 参考常见问题章节
4. 联系技术支持

---

**祝你部署顺利！** 🚀
