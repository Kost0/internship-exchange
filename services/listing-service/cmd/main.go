package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	listingpb "github.com/Kost0/internship-exchange/proto/listing"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/config"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/handler"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/repository"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/service"
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

	listingRepo := repository.NewListingRepository(pool)
	companyRepo := repository.NewCompanyRepository(pool)
	listingSvc := service.NewListingService(listingRepo, companyRepo)
	listingHandler := handler.NewListingHandler(listingSvc)

	grpcServer := grpc.NewServer()
	listingpb.RegisterListingServiceServer(grpcServer, listingHandler)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("listing-service gRPC starting on %s", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migration, err := os.ReadFile("migrations/001_create_tables.up.sql")
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, string(migration))

	return err
}
