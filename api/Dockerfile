# =====================
# 1. Build stage
# =====================
FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY api ./api
COPY common ./common

RUN go build -o /api-server ./api

# =====================
# 2. Runtime stage
# =====================
FROM alpine:latest

WORKDIR /app

# Install ca-certificates, curl, and bash for Atlas installation
RUN apk --no-cache add ca-certificates curl bash

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

# Copy the compiled binary
COPY --from=builder /api-server /app/api-server

# Copy migration files
COPY Database/migrations /app/Database/migrations

EXPOSE 8000

CMD ["/app/api-server"]
