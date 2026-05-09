package main

import (
	"log"
	"net/http"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/config"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/router"
)

func main() {
	cfg := config.Load()

	r := router.New(cfg)

	log.Printf("API Gateway starting on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
