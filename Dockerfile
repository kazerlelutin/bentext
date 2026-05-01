FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# go.sum may not exist if the module has no dependencies yet
COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Pinned runtime image (avoid :latest for reproducibility and supply-chain clarity)
FROM alpine:3.21

RUN apk --no-cache add ca-certificates wget

RUN addgroup -g 1001 -S bentext && \
  adduser -u 1001 -S bentext -G bentext

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/ingredient-sprites.bentext ./ingredient-sprites.bentext
COPY --from=builder /app/recipes ./recipes
COPY --from=builder /app/public ./public
COPY --from=builder /app/docker-entrypoint.sh ./docker-entrypoint.sh

RUN chown -R bentext:bentext /app && \
  chmod +x ./docker-entrypoint.sh && \
  chmod +x ./main

USER bentext

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./docker-entrypoint.sh"]
