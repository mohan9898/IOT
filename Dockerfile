# Build stage for backend
FROM golang:1.21-alpine AS backend-builder

# Install CGO dependencies
RUN apk add --no-cache gcc musl-dev git make

WORKDIR /app

# Copy go mod and sum first (caching layer)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build with verbose logging and debugging
RUN ls -la && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o iot-manager main.go

# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Install dependencies first (caching layer)
COPY frontend/package*.json ./
RUN npm ci

# Copy frontend source and build
COPY frontend/ ./
RUN npm run build

# Final production stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs tzdata

WORKDIR /root/

# Copy backend binary from builder stage
COPY --from=backend-builder /app/iot-manager .

# Copy frontend build
COPY --from=frontend-builder /dist/ ./dist/

# Create and prepare data directory
RUN mkdir -p /root/data && \
    chmod 755 /root/data

# Expose port
EXPOSE 6116

# Default environment variables (non-sensitive)
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=6116
ENV HTTPS_ENABLE=false
ENV DB_PATH=/root/data/iot.db
ENV DB_BACKUP_ENABLE=false
ENV DB_BACKUP_PATH=/root/data/backup

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:6116/health || exit 1

# Start server
CMD ["./iot-manager"]
