
# IoT Device Manager

A simple, production-ready IoT device management system built with Go + Vue 3.

## 🚀 Quick Start

### Requirements
- Docker and Docker Compose

### Deploy with 1 command
```bash
# 1. Clone and enter repo
git clone https://github.com/mohan9898/IOT.git
cd IOT

# 2. Copy and edit environment variables
cp .env.example .env
# Edit .env - set JWT_SECRET and MQTT credentials

# 3. Start the service
docker-compose up -d

# 4. Access
http://localhost:6116
```

## 📦 Quick Deploy Script
```bash
chmod +x deploy.sh
./deploy.sh
```

## 🔧 Configuration
Copy `.env.example` to `.env` and edit:
- `JWT_SECRET` - **REQUIRED** Secure secret for JWT tokens
- `MQTT_USERNAME` - MQTT broker username
- `MQTT_PASSWORD` - MQTT broker password
- Other settings are optional with defaults

## 🐳 Using the pre-built image
The image is automatically built and pushed to GitHub Container Registry:
```
ghcr.io/mohan9898/iot:latest
```

## 📂 Project Structure
```
IOT/
├── backend/          # Go backend API
├── frontend/         # Vue 3 frontend source
├── dist/            # Built frontend (included)
├── Dockerfile       # Docker build file
├── docker-compose.yml
└── .env.example
```

## 🛠️ Development
If you want to build from source:

```bash
# Backend
cd backend
go mod download
go run main.go

# Frontend
cd frontend
npm install
npm run dev
```

## 🔒 Security
- Runs as non-root user inside container
- JWT authentication
- No sensitive data committed
- Secure defaults

## 📝 License
MIT
