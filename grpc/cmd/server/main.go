package main

import (
	"context"
	"flag"
	noteAPI "github.com/igorakimy/bigtech_microservices/internal/api/note"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/igorakimy/bigtech_microservices/internal/config/env"
	nr "github.com/igorakimy/bigtech_microservices/internal/repository/note"
	noteService "github.com/igorakimy/bigtech_microservices/internal/service/note"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Read env variables
	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPostgresConfig()
	if err != nil {
		log.Fatalf("failed to get postgres config: %v", err)
	}

	// Create listener
	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create database pool
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	noteRepo := nr.NewPostgresRepository(pool)
	noteServ := noteService.NewService(noteRepo)

	// Register server
	srv := grpc.NewServer()
	reflection.Register(srv)
	desc.RegisterNoteV1Server(srv, noteAPI.NewServerImplementation(noteServ))

	log.Printf("server listening at: %v", lis.Addr())

	// Serve server
	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
