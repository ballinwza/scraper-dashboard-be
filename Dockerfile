# --- Base Stage ---
FROM golang:1.25.0-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# --- Development Stage (ใช้สำหรับ dev) ---
FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
CMD ["air", "-c", ".air.toml"]

# --- Builder Stage (คอมไพล์เพื่อ prod) ---
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# --- Production Runner Stage ---
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]