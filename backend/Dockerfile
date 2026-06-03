# syntax=docker/dockerfile:1.6
FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/cli ./cli

FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /out/cli /app/cli
COPY --from=builder /src/migrations /app/migrations
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

RUN chmod +x /app/docker-entrypoint.sh

USER app

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["--config=/app/config.yml"]
