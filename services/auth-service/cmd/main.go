package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	authpb "github.com/Kost0/internship-exchange/proto/auth"
	"github.com/Kost0/internship-exchange/services/auth-service/internal/config"
	"github.com/Kost0/internship-exchange/services/auth-service/internal/handler"
	"github.com/Kost0/internship-exchange/services/auth-service/internal/repository"
	"github.com/Kost0/internship-exchange/services/auth-service/internal/service"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("database connected")

	if err := runMigrations(context.Background(), pool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied")

	accessTTL, _ := time.ParseDuration(cfg.AccessTokenTTL)
	refreshTTL, _ := time.ParseDuration(cfg.RefreshTokenTTL)

	userRepo := repository.NewUserRepository(pool)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, accessTTL, refreshTTL)
	authHandler := handler.NewAuthHandler(authSvc)

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("auth-service gRPC starting on %s", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migration, err := os.ReadFile("migrations/001_create_users.up.sql")
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, string(migration))
	return err
}
