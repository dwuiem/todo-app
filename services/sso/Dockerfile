FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ../.. .

RUN go build -o sso ./cmd/main.go

EXPOSE 8081

CMD ["./sso"]
