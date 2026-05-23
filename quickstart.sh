#!/bin/bash
# IoT Manager 快速启动脚本
# 适用于快速部署和测试

set -e

echo "=========================================="
echo "   IoT Manager 快速部署脚本"
echo "=========================================="
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    echo "访问 https://docs.docker.com/get-docker/ 了解更多"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 创建必要目录
mkdir -p data
mkdir -p data/backup

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "📝 创建 .env 文件..."
    cp .env.example .env
    
    # 生成随机的 JWT_SECRET
    if command -v openssl &> /dev/null; then
        RANDOM_SECRET=$(openssl rand -hex 32)
        sed -i.bak "s/your-very-secure-secret-key-change-in-production/$RANDOM_SECRET/" .env
        rm -f .env.bak
        echo "✅ 已生成随机 JWT 密钥"
    fi
else
    echo "✅ .env 文件已存在"
fi

echo ""
echo "选择部署方式："
echo "1) 完整部署（IoT Manager + Prometheus + Grafana）- 推荐"
echo "2) 仅部署 IoT Manager（轻量版）"
echo "3) 本地开发环境"
echo ""
read -p "请输入选项 (1-3): " choice

case $choice in
    1)
        echo ""
        echo "🚀 开始完整部署..."
        docker-compose up -d
        echo ""
        echo "✅ 完整部署完成！"
        echo ""
        echo "📱 访问地址："
        echo "   IoT Manager:  http://localhost:6116"
        echo "   Prometheus:   http://localhost:9090"
        echo "   Grafana:      http://localhost:3000 (admin/admin)"
        echo ""
        ;;
        
    2)
        echo ""
        echo "🚀 开始轻量部署（仅 IoT Manager）..."
        docker-compose -f docker-compose-core.yml up -d
        echo ""
        echo "✅ 轻量部署完成！"
        echo ""
        echo "📱 访问地址："
        echo "   IoT Manager:  http://localhost:6116"
        echo ""
        ;;
        
    3)
        echo ""
        echo "🔧 准备本地开发环境..."
        
        # 检查 Go
        if ! command -v go &> /dev/null; then
            echo "❌ Go 未安装，请先安装 Go"
            exit 1
        fi
        
        # 检查 Node.js
        if ! command -v node &> /dev/null; then
            echo "❌ Node.js 未安装，请先安装 Node.js"
            exit 1
        fi
        
        echo "✅ 依赖检查通过"
        echo ""
        echo "📝 下一步操作："
        echo "   cd backend && go mod download && go run main.go"
        echo "   cd frontend && npm install && npm run dev"
        echo ""
        ;;
        
    *)
        echo "❌ 无效选项"
        exit 1
        ;;
esac

echo ""
echo "💡 提示："
echo "   查看日志: docker-compose logs -f"
echo "   停止服务: docker-compose stop"
echo "   更多信息: 查看 DEPLOYMENT.md"
echo ""
echo "=========================================="
