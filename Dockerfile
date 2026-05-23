
# IoT Manager - Core Build (Simplified)
# Uses existing dist/ folder, no need to build frontend!

# Stage 1: Build Go backend only
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

# Copy go modules first
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source and build
COPY backend/ ./
RUN ls -la
RUN CGO_ENABLED=1 GOOS=linux go build -v -o iot-manager .

# Stage 2: Final image - using existing dist/ from repo!
FROM alpine:3.19

# Install dependencies
RUN apk add --no-cache ca-certificates sqlite-libs curl

# Set up non-root user
RUN addgroup -g 1000 iot && adduser -D -u 1000 -G iot -h /app iot
WORKDIR /app

# Copy backend binary
COPY --from=builder /app/iot-manager .

# Copy existing dist folder directly from repo!
COPY dist/ ./dist/

# Prepare data directory
RUN mkdir -p /app/data && chown -R iot:iot /app

# Switch user
USER iot

# Expose port
EXPOSE 6116

# Environment
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=6116
ENV DB_PATH=/app/data/iot.db
ENV GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=45s --retries=3 \
  CMD curl -f http://localhost:6116/health || exit 1

# Run
CMD ["/app/iot-manager"]
