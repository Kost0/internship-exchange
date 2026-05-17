package config

import "os"

type Config struct {
	GRPCAddr    string
	DatabaseURL string
	RabbitMQURL string
}

func Load() *Config {
	return &Config{
		GRPCAddr:    getEnv("GRPC_ADDR", ":50054"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/application_db?sslmode=disable"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	
	return fallback
}
