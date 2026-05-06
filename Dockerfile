# ── Build stage ────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o waweb-v2 .

# ── Runtime stage ───────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /app/waweb-v2 .

# Copy runtime assets
COPY --from=builder /app/views ./views
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Create writable directories
RUN mkdir -p uploads published_sites

EXPOSE 3000

CMD ["./waweb-v2"]
