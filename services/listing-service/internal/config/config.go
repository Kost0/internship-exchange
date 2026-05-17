package config

import "os"

type Config struct {
	GRPCAddr    string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		GRPCAddr:    getEnv("GRPC_ADDR", ":50053"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/listing_db?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
