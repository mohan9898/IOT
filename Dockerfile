FROM golang:1.21-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY backend/ ./

RUN go env && echo "---files---" && ls -la && echo "---download---"

RUN go mod download && go mod verify

RUN echo "---build---" && CGO_ENABLED=1 GOOS=linux go build -v -o iot-manager .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 curl && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 1000 iot && useradd -u 1000 -g iot -d /app -s /usr/sbin/nologin iot

WORKDIR /app

COPY --from=builder /build/iot-manager .
COPY dist/ ./dist/

RUN mkdir -p /app/data && chown -R iot:iot /app

USER iot

EXPOSE 6116

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=6116
ENV DB_PATH=/app/data/iot.db
ENV GIN_MODE=release

HEALTHCHECK --interval=30s --timeout=10s --start-period=45s --retries=3 \
  CMD curl -f http://localhost:6116/health || exit 1

CMD ["/app/iot-manager"]