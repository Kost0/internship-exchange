FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o profile-service ./services/profile-service/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/profile-service .
COPY --from=builder /app/services/profile-service/migrations ./migrations

EXPOSE 50052

CMD ["./profile-service"]