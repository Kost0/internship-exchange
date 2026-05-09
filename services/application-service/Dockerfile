FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/auth-service .

COPY --from=builder /app/migrations ./migrations

EXPOSE 50051

CMD ["./auth-service"]