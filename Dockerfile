# ---------- builder ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary
COPY --from=builder /app/app /app/app

# Copy migrations (read-only)
COPY --from=builder /app/db/migrations /app/db/migrations

# Railway & docker-compose friendly
EXPOSE 8000

CMD ["/app/app"]
