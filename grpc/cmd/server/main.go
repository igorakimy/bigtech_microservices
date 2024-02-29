package main

import (
	"context"
	"database/sql"
	"flag"
	sq "github.com/Masterminds/squirrel"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/igorakimy/bigtech_microservices/internal/config/env"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedNoteV1Server
	pool *pgxpool.Pool
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	buildGetNote := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := buildGetNote.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var id int64
	var title, body string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).
		Scan(&id, &title, &body, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to build sql: %v", err)
	}

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: id,
			Info: &desc.NoteInfo{
				Title:   title,
				Content: body,
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updatedAtTime,
		},
	}, nil
}

func (s *server) List(ctx context.Context, req *desc.ListRequest) (*desc.ListResponse, error) {
	buildListNotes := sq.Select("id", "title", "body", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("note")

	query, args, err := buildListNotes.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to get list of notes: %v", err)
	}

	var n = struct {
		id          int64
		title, body string
		createdAt   time.Time
		updatedAt   sql.NullTime
	}{}

	var notes []*desc.Note

	for rows.Next() {
		if err = rows.Scan(&n.id, &n.title, &n.body, &n.createdAt, &n.updatedAt); err != nil {
			log.Fatalf("failed to scan data: %v", err)
		}

		var note = &desc.Note{
			Id: n.id,
			Info: &desc.NoteInfo{
				Title:   n.title,
				Content: n.body,
			},
			CreatedAt: timestamppb.New(n.createdAt),
			UpdatedAt: timestamppb.New(n.updatedAt.Time),
		}

		notes = append(notes, note)
	}

	return &desc.ListResponse{Notes: notes}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	buildCreateNote := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body", "created_at").
		Values(req.GetInfo().GetTitle(), req.GetInfo().GetContent(), time.Now()).
		Suffix("RETURNING id")

	query, args, err := buildCreateNote.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var noteID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&noteID)
	if err != nil {
		log.Fatalf("failed to create note: %v", err)
	}

	return &desc.CreateResponse{
		Id: noteID,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	buildUpdateNote := sq.Update("note").
		PlaceholderFormat(sq.Dollar).
		Set("title", req.GetInfo().GetTitle().GetValue()).
		Set("body", req.GetInfo().GetContent().GetValue()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := buildUpdateNote.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	buildDeleteNote := sq.Delete("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := buildDeleteNote.ToSql()
	if err != nil {
		log.Fatalf("failed to build delete query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete note: %v", err)
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
	desc.RegisterNoteV1Server(srv, &server{pool: pool})

	log.Printf("server listening at: %v", lis.Addr())

	// Serve server
	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
