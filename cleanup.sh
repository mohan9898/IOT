#!/bin/bash
# IoT Manager 停止和清理脚本

set -e

echo "=========================================="
echo "   IoT Manager 停止和清理"
echo "=========================================="
echo ""

# 检查是否有容器在运行
if [ "$(docker-compose ps -q)" ]; then
    echo "选择操作："
    echo "1) 仅停止服务"
    echo "2) 停止并删除容器（保留数据）"
    echo "3) 完全清除（删除所有数据）"
    echo ""
    read -p "请输入选项 (1-3): " choice
    
    case $choice in
        1)
            echo ""
            echo "⏹️  正在停止服务..."
            docker-compose stop
            echo "✅ 服务已停止"
            ;;
            
        2)
            echo ""
            echo "🗑️  正在停止并删除容器..."
            docker-compose down
            echo "✅ 容器已删除，数据已保留"
            ;;
            
        3)
            echo ""
            echo "⚠️  警告：这将删除所有数据和配置！"
            read -p "确认继续? (yes/no): " confirm
            
            if [ "$confirm" = "yes" ]; then
                echo ""
                echo "🗑️  正在完全清除..."
                docker-compose down -v
                docker rmi iot-manager 2>/dev/null || true
                echo "✅ 已完全清除"
            else
                echo "操作已取消"
            fi
            ;;
            
        *)
            echo "❌ 无效选项"
            exit 1
            ;;
    esac
else
    echo "ℹ️  没有服务在运行"
fi

echo ""
echo "=========================================="
