package config

import "os"

type Config struct {
	GRPCAddr           string
	DatabaseURL        string
	MinioAddr          string
	MinioUser          string
	MinioPass          string
	MinioBucket        string
	ListingServiceAddr string
}

func Load() *Config {
	return &Config{
		GRPCAddr:           getEnv("GRPC_ADDR", ":50052"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/profile_db?sslmode=disable"),
		MinioAddr:          getEnv("MINIO_ADDR", "localhost:9000"),
		MinioUser:          getEnv("MINIO_USER", "minioadmin"),
		MinioPass:          getEnv("MINIO_PASS", "minioadmin"),
		MinioBucket:        getEnv("MINIO_BUCKET", "internship-exchange"),
		ListingServiceAddr: getEnv("LISTING_SERVICE_ADDR", "localhost:50053"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
