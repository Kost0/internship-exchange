FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth-service ./services/auth-service/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/auth-service .
COPY --from=builder /app/services/auth-service/migrations ./migrations

EXPOSE 50051

CMD ["./auth-service"]