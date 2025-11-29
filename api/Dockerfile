FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY common ./common
COPY api ./api

RUN go install github.com/cosmtrek/air@v1.42.0

CMD ["air", "--build.cmd", "go build -o ./tmp/api ./api", "--build.bin", "./tmp/api"]
