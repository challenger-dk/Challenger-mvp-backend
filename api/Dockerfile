# =====================
# 1. Build stage
# =====================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy all code
COPY api ./api
COPY common ./common

# Build production binary
RUN go build -o /api-server ./api

# =====================
# 2. Runtime stage
# =====================
FROM alpine

WORKDIR /app

COPY --from=builder /api-server .

EXPOSE 8000

CMD ["./api-server"]
