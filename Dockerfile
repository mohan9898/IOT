
# Multi-stage build for IoT Manager
# Stage 1: Build backend
FROM golang:1.21-alpine AS go-builder
RUN apk add --no-cache gcc musl-dev git
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN ls -la
RUN CGO_ENABLED=1 GOOS=linux go build -v -o iot-manager main.go

# Stage 2: Build frontend
FROM node:20-alpine AS node-builder
WORKDIR /src
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 3: Production image
FROM alpine:3.19
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app

# Copy binaries and assets
COPY --from=go-builder /src/iot-manager /app/iot-manager
COPY --from=node-builder /src/dist /app/dist

# Prepare data directory
RUN mkdir -p /app/data && chmod 755 /app/data

# Expose port
EXPOSE 6116

# Environment variables
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=6116
ENV DB_PATH=/app/data/iot.db

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:6116/health || exit 1

# Run
CMD ["/app/iot-manager"]
