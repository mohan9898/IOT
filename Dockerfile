# Build stage for backend
FROM golang:1.21-alpine AS backend-builder

# Install CGO dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -o iot-manager main.go

# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci

# Copy frontend source
COPY frontend/ ./
RUN npm run build

# Final production stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy backend binary
COPY --from=backend-builder /app/iot-manager .

# Copy frontend build (from /dist since WORKDIR=/app in frontend-builder)
COPY --from=frontend-builder /dist/ ./dist/

# Create data directory
RUN mkdir -p /root/data

# Expose port
EXPOSE 6116

# Default environment variables
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=6116
ENV HTTPS_ENABLE=false
ENV DB_PATH=/root/data/iot.db
ENV DB_BACKUP_ENABLE=false
ENV DB_BACKUP_PATH=/root/data/backup
ENV JWT_SECRET=change-this-in-production

# Start server
CMD ["./iot-manager"]
