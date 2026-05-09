package config

import "os"

type Config struct {
	GRPCAddr        string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  string
	RefreshTokenTTL string
}

func Load() *Config {
	return &Config{
		GRPCAddr:        getEnv("GRPC_ADDR", ":50051"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "secret"),
		AccessTokenTTL:  getEnv("ACCESS_TOKEN_TTL", "15m"),
		RefreshTokenTTL: getEnv("REFRESH_TOKEN_TTL", "720h"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
