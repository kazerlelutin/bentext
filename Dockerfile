FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

RUN addgroup -g 1001 -S bentext && \
    adduser -u 1001 -S bentext -G bentext

WORKDIR /app


COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

COPY --from=builder /app/static ./static

COPY --from=builder /app/recipes ./recipes
COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/docker-entrypoint.sh ./docker-entrypoint.sh

RUN touch .env

RUN chown -R bentext:bentext /app && \
    chmod +x ./docker-entrypoint.sh && \
    chmod +x ./main && \
    chmod +x ./migrate

USER bentext

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["./docker-entrypoint.sh"]
