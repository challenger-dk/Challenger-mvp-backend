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
FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /api-server /app/api-server

EXPOSE 8000

CMD ["/app/api-server"]
