FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o listing-service ./services/listing-service/cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/listing-service .
COPY --from=builder /app/services/listing-service/migrations ./migrations

EXPOSE 50053

CMD ["./listing-service"]