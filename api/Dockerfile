# =====================
# 1. Build stage
# =====================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Ensure static build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY api ./api
COPY common ./common

# Build production binary
RUN go build -o /api-server ./api

# =====================
# 2. Runtime stage (distroless)
# =====================
FROM gcr.io/distroless/static

COPY --from=builder /api-server /api-server

# Cloud Run will set PORT env variable
EXPOSE 8000

CMD ["/api-server"]
