FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/migrator ./cmd/migrator
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/app ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/migrator /app/migrator
COPY --from=builder /app/app /app/app

COPY internal/storage/clickhouse/migrations ./internal/storage/clickhouse/migrations
COPY .env .

LABEL authors="recrusion"

CMD ["/app/app"]