FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN go install github.com/cosmtrek/air@v1.42.0

COPY . .

EXPOSE 8080
CMD ["air"]
