package main

import (
	"context"
	"errors"
	"flag"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/igorakimy/bigtech_microservices/internal/config/env"
	"github.com/igorakimy/bigtech_microservices/internal/converter"
	nr "github.com/igorakimy/bigtech_microservices/internal/repository/note"
	"github.com/igorakimy/bigtech_microservices/internal/service"
	noteService "github.com/igorakimy/bigtech_microservices/internal/service/note"
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
	noteService service.NoteService
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	note, err := s.noteService.Get(ctx, req.GetId())

	if err != nil {
		return nil, err
	}

	if note == nil {
		return nil, errors.New("note not found")
	}

	return &desc.GetResponse{
		Note: converter.ToNoteFromService(note),
	}, nil
}

func (s *server) List(ctx context.Context, _ *desc.ListRequest) (*desc.ListResponse, error) {
	notes, err := s.noteService.List(ctx)
	if err != nil {
		return nil, err
	}
	return &desc.ListResponse{
		Notes: converter.ToNotesFromService(notes),
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	noteID, err := s.noteService.Create(
		ctx,
		converter.ToNoteInfoFromDesc(req.GetInfo()),
	)
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: noteID}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := s.noteService.Update(
		ctx,
		req.GetId(),
		converter.ToUpdateNoteInfoFromDesc(req.GetInfo()),
	)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.noteService.Delete(ctx, req.GetId()); err != nil {
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

	noteRepo := nr.NewPostgresRepository(pool)
	noteServ := noteService.NewService(noteRepo)

	// Register server
	srv := grpc.NewServer()
	reflection.Register(srv)
	desc.RegisterNoteV1Server(srv, &server{
		noteService: noteServ,
	})

	log.Printf("server listening at: %v", lis.Addr())

	// Serve server
	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
