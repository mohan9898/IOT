#!/bin/bash
# IoT Manager 系统状态检查脚本

echo "=========================================="
echo "   IoT Manager 系统状态检查"
echo "=========================================="
echo ""

# 检查 Docker
echo "📦 Docker 状态："
if command -v docker &> /dev/null; then
    echo "   ✅ Docker 已安装: $(docker --version)"
    if docker info &> /dev/null; then
        echo "   ✅ Docker 正在运行"
    else
        echo "   ❌ Docker 未运行"
    fi
else
    echo "   ❌ Docker 未安装"
fi

echo ""

# 检查 Docker Compose
echo "📦 Docker Compose 状态："
if command -v docker-compose &> /dev/null; then
    echo "   ✅ Docker Compose 已安装: $(docker-compose --version)"
else
    echo "   ❌ Docker Compose 未安装"
fi

echo ""

# 检查项目文件
echo "📁 项目文件："
for file in docker-compose.yml DEPLOYMENT.md README.md; do
    if [ -f "$file" ]; then
        echo "   ✅ $file 存在"
    else
        echo "   ❌ $file 不存在"
    fi
done

echo ""

# 检查容器状态
echo "🐳 容器状态："
if [ -f docker-compose.yml ]; then
    containers=$(docker-compose ps --services 2>/dev/null | wc -l)
    if [ "$containers" -gt 0 ]; then
        docker-compose ps
        echo ""
        echo "   运行中的容器："
        docker-compose ps --format json | jq -r '.[] | select(.State == "running") | "   ✅ \(.Service) - 运行中"'
        echo ""
        echo "   停止的容器："
        docker-compose ps --format json | jq -r '.[] | select(.State != "running") | "   ❌ \(.Service) - \(.State)"'
    else
        echo "   ℹ️  没有容器在运行"
    fi
else
    echo "   ⚠️  docker-compose.yml 不存在"
fi

echo ""

# 检查端口占用
echo "🔌 端口状态："
ports=(6116 3000 9090 9093)
for port in "${ports[@]}"; do
    if command -v netstat &> /dev/null; then
        if netstat -tlnp 2>/dev/null | grep -q ":$port "; then
            echo "   ✅ 端口 $port 被占用"
        else
            echo "   ℹ️  端口 $port 可用"
        fi
    elif command -v ss &> /dev/null; then
        if ss -tlnp 2>/dev/null | grep -q ":$port "; then
            echo "   ✅ 端口 $port 被占用"
        else
            echo "   ℹ️  端口 $port 可用"
        fi
    fi
done

echo ""

# 检查磁盘空间
echo "💾 磁盘空间："
if command -v df &> /dev/null; then
    df -h . | tail -n 1 | awk '{printf "   已用: %s/%s (可用: %s)\n", $3, $2, $4}'
fi

echo ""

# 检查监控目录
echo "📊 监控配置："
if [ -d "monitoring" ]; then
    echo "   ✅ monitoring 目录存在"
    for dir in prometheus grafana; do
        if [ -d "monitoring/$dir" ]; then
            echo "   ✅ monitoring/$dir 存在"
        else
            echo "   ❌ monitoring/$dir 不存在"
        fi
    done
else
    echo "   ❌ monitoring 目录不存在"
fi

echo ""
echo "=========================================="
echo ""
echo "💡 下一步："
echo "   部署: ./quickstart.sh"
echo "   清理: ./cleanup.sh"
echo "   文档: 查看 DEPLOYMENT.md"
echo ""
