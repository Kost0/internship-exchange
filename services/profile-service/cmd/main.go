package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/Kost0/internship-exchange/services/profile-service/internal/clients"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	profilepb "github.com/Kost0/internship-exchange/proto/profile"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/config"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/handler"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/repository"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/service"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/storage"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("database connected")

	if err = runMigrations(context.Background(), pool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied")

	minioStorage, err := storage.NewMinioStorage(cfg.MinioAddr, cfg.MinioUser, cfg.MinioPass)
	if err != nil {
		log.Fatalf("failed to connect to minio: %v", err)
	}

	if err = minioStorage.EnsureBuckets(context.Background()); err != nil {
		log.Fatalf("failed to ensure buckets: %v", err)
	}
	log.Println("minio connected")

	listingClient, err := clients.NewListingClient(cfg.ListingServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to listing-service: %v", err)
	}

	studentRepo := repository.NewStudentRepository(pool)
	companyRepo := repository.NewCompanyRepository(pool)
	profileSvc := service.NewProfileService(studentRepo, companyRepo, minioStorage, listingClient)
	profileHandler := handler.NewProfileHandler(profileSvc)

	grpcServer := grpc.NewServer()
	profilepb.RegisterProfileServiceServer(grpcServer, profileHandler)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("profile-service gRPC starting on %s", cfg.GRPCAddr)
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
