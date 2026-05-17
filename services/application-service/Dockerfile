FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o application-service ./services/application-service/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/application-service .
COPY --from=builder /app/services/application-service/migrations ./migrations

EXPOSE 50054

CMD ["./application-service"]