FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY common ./common
COPY chat ./chat

RUN go install github.com/cosmtrek/air@v1.42.0

CMD ["air", "--build.cmd", "go build -o ./tmp/chat ./chat", "--build.bin", "./tmp/chat"]
