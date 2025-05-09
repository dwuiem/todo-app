FROM golang:1.24.2-alpine

ENV CONFIG_PATH=/app/config/local.yaml

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o migrate ./cmd/migrate/main.go
RUN go build -o todo ./cmd/app/main.go

EXPOSE 8080

CMD ./migrate && ./todo

