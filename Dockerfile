FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev git sqlite-dev

WORKDIR /build

COPY backend/ ./

RUN go env && echo "---" && ls -la

RUN go mod download && go mod verify

RUN CGO_ENABLED=1 GOOS=linux go build -mod=mod -v -o iot-manager .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates sqlite-libs curl && \
    addgroup -g 1000 iot && adduser -D -u 1000 -G iot -h /app iot

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