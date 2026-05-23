
#!/bin/bash
# IoT Manager - Quick Deploy Script
set -e

echo "=========================================="
echo "  🚀 IoT Manager Quick Deploy"
echo "=========================================="
echo ""

# Check Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is required but not installed"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is required but not installed"
    exit 1
fi

echo "✅ Docker check passed"
echo ""

# Check .env file
if [ ! -f .env ]; then
    echo "⚠️  .env file not found, creating from .env.example..."
    cp .env.example .env
    echo ""
    echo "⚠️  Please edit .env and set JWT_SECRET, MQTT_USERNAME, MQTT_PASSWORD"
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

# Choose compose command
COMPOSE_CMD="docker-compose"
if ! command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker compose"
fi

# Pull and start
echo ""
echo "📦 Pulling latest image..."
$COMPOSE_CMD pull

echo ""
echo "🚀 Starting service..."
$COMPOSE_CMD up -d

echo ""
echo "=========================================="
echo "  ✅ Deployment complete!"
echo "=========================================="
echo ""
echo "📡 Access IoT Manager at: http://localhost:6116"
echo ""
echo "📝 View logs: $COMPOSE_CMD logs -f"
echo "🔄 Update: $COMPOSE_CMD pull && $COMPOSE_CMD up -d"
echo ""
