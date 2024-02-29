package main

import (
	"context"
	"flag"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/igorakimy/bigtech_microservices/internal/config/env"
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	nr "github.com/igorakimy/bigtech_microservices/internal/repository/note"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedNoteV1Server
	noteRepo repository.NoteRepository
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	note, err := s.noteRepo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &desc.GetResponse{Note: note}, nil
}

func (s *server) List(ctx context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {
	notes, err := s.noteRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return &desc.ListResponse{Notes: notes}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	noteID, err := s.noteRepo.Create(ctx, req.Info)
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: noteID}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if err := s.noteRepo.Update(ctx, req.GetId(), req.GetInfo()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.noteRepo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
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

	// Register server
	srv := grpc.NewServer()
	reflection.Register(srv)
	desc.RegisterNoteV1Server(srv, &server{
		noteRepo: nr.NewPostgresRepository(pool),
	})

	log.Printf("server listening at: %v", lis.Addr())

	// Serve server
	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
