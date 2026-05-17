package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kost0/internship-exchange/services/notification-service/internal/config"
	"github.com/Kost0/internship-exchange/services/notification-service/internal/consumer"
	"github.com/Kost0/internship-exchange/services/notification-service/internal/mailer"
)

func main() {
	cfg := config.Load()

	m := mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.FromEmail)

	c, err := consumer.New(cfg.RabbitMQURL, m)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer c.Close()
	log.Println("notification-service started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("shutting down...")
		cancel()
	}()

	if err := c.Start(ctx); err != nil {
		log.Fatalf("consumer error: %v", err)
	}
}
