package config

import (
	"os"
)

type Config struct {
	HTTPAddr           string
	JWTSecret          string
	RedisAddr          string
	AuthServiceAddr    string
	ProfileServiceAddr string
	ListingServiceAddr string
	AppServiceAddr     string
}

func Load() *Config {
	return &Config{
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		AuthServiceAddr:    getEnv("AUTH_SERVICE_ADDR", "localhost:50051"),
		ProfileServiceAddr: getEnv("PROFILE_SERVICE_ADDR", "localhost:50052"),
		ListingServiceAddr: getEnv("LISTING_SERVICE_ADDR", "localhost:50053"),
		AppServiceAddr:     getEnv("APP_SERVICE_ADDR", "localhost:50054"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
