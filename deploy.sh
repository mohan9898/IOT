#!/bin/bash
set -e

# IoT Manager 快速部署脚本
# 使用 GitHub Container Registry 镜像

echo "=========================================="
echo "  🚀 IoT Manager 快速部署脚本"
echo "=========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查 Docker
echo -e "${BLUE}检查 Docker 和 Docker Compose...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误：Docker 未安装${NC}"
    echo "请先安装 Docker：https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${RED}错误：Docker Compose 未安装${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker 检查通过${NC}"
echo ""

# 选择部署模式
echo "请选择部署模式："
echo "1) 核心版本（仅 IoT Manager）- 推荐"
echo "2) 完整版本（含 Prometheus + Grafana + Alertmanager）"
read -p "请输入选项 (1-2，默认 1): " mode
mode=${mode:-1}

# 检查 .env 文件
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env 文件不存在，从 .env.example 创建...${NC}"
    cp .env.example .env
    echo ""
    echo -e "${YELLOW}请配置 .env 文件中的重要参数：${NC}"
    echo "  - JWT_SECRET"
    echo "  - MQTT_USERNAME"
    echo "  - MQTT_PASSWORD"
    echo ""
    
    # 生成随机 JWT_SECRET
    if command -v openssl &> /dev/null; then
        RANDOM_SECRET=$(openssl rand -base64 32)
        sed -i.bak "s/^JWT_SECRET=.*/JWT_SECRET=$RANDOM_SECRET/" .env
        rm -f .env.bak
        echo -e "${GREEN}✓ 已自动生成随机 JWT_SECRET${NC}"
    fi
    
    echo ""
    read -p "是否继续部署？(y/n): " continue
    if [ "$continue" != "y" ] && [ "$continue" != "Y" ]; then
        echo "部署已取消。"
        exit 0
    fi
fi

# 选择 docker-compose 命令
COMPOSE_CMD="docker-compose"
if ! command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker compose"
fi

# 选择配置文件
if [ "$mode" = "2" ]; then
    COMPOSE_FILE="docker-compose-ghcr-full.yml"
    echo -e "${BLUE}部署完整版本...${NC}"
else
    COMPOSE_FILE="docker-compose-ghcr.yml"
    echo -e "${BLUE}部署核心版本...${NC}"
fi

# 拉取最新镜像
echo ""
echo -e "${BLUE}拉取最新 Docker 镜像...${NC}"
$COMPOSE_CMD -f $COMPOSE_FILE pull

# 启动服务
echo ""
echo -e "${BLUE}启动服务...${NC}"
$COMPOSE_CMD -f $COMPOSE_FILE up -d

# 等待服务启动
echo ""
echo -e "${BLUE}等待服务启动...${NC}"
sleep 10

# 检查状态
echo ""
echo -e "${BLUE}检查服务状态...${NC}"
$COMPOSE_CMD -f $COMPOSE_FILE ps

# 显示完成信息
echo ""
echo "=========================================="
echo -e "${GREEN}  🎉 部署成功！${NC}"
echo "=========================================="
echo ""
echo "📦 服务地址："
echo "   IoT Manager: http://localhost:6116"
if [ "$mode" = "2" ]; then
    echo "   Grafana:     http://localhost:3000"
    echo "   Prometheus:  http://localhost:9090"
fi
echo ""
echo "📝 管理命令："
echo "   查看日志:        $COMPOSE_CMD -f $COMPOSE_FILE logs -f"
echo "   停止服务:        $COMPOSE_CMD -f $COMPOSE_FILE down"
echo "   更新镜像:        $COMPOSE_CMD -f $COMPOSE_FILE pull && $COMPOSE_CMD -f $COMPOSE_FILE up -d"
echo ""
echo "🔧 配置文件："
echo "   环境变量:        .env"
echo "   Docker Compose:  $COMPOSE_FILE"
echo ""
