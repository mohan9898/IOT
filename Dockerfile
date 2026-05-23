
# ==========================================
# IoT Manager - Production Dockerfile
# ==========================================
# Stage 1: Build Go Backend
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Set working directory
WORKDIR /build

# Copy go mod files first for caching
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Verify file structure
RUN ls -la && go version

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -v -o iot-manager .

# ==========================================
# Stage 2: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /build
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ==========================================
# Stage 3: Production Image
FROM alpine:3.19

# Set maintainer label
LABEL maintainer="IoT Manager Team"
LABEL description="IoT Device Management System"

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs tzdata curl

# Create non-root user for security
RUN addgroup -g 1000 iot && \
    adduser -D -u 1000 -G iot -h /app -s /sbin/nologin iot

# Set working directory
WORKDIR /app

# Copy artifacts from build stages
COPY --from=builder --chown=iot:iot /build/iot-manager ./
COPY --from=frontend-builder --chown=iot:iot /build/dist ./dist

# Create and prepare data directory
RUN mkdir -p /app/data && \
    chown -R iot:iot /app

# Switch to non-root user
USER iot

# Expose port
EXPOSE 6116

# Environment variables
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=6116
ENV DB_PATH=/app/data/iot.db
ENV GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=45s --retries=3 \
  CMD curl -f http://localhost:6116/health || exit 1

# Start the application
CMD ["/app/iot-manager"]
