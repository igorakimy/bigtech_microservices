package app

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/api/note"
	"github.com/igorakimy/bigtech_microservices/internal/closer"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/igorakimy/bigtech_microservices/internal/config/env"
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	noteRepository "github.com/igorakimy/bigtech_microservices/internal/repository/note"
	"github.com/igorakimy/bigtech_microservices/internal/service"
	noteService "github.com/igorakimy/bigtech_microservices/internal/service/note"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type serviceProvider struct {
	postgresConfig config.PostgresConfig
	grpcConfig     config.GRPCConfig

	postgresPool *pgxpool.Pool
	noteRepo     repository.NoteRepository

	noteService service.NoteService

	noteImpl *note.ServerImplementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PostgresConfig() config.PostgresConfig {
	if s.postgresConfig == nil {
		pgConfig, err := env.NewPostgresConfig()
		if err != nil {
			log.Fatalf("failed to get postgres config: %s", err.Error())
		}
		s.postgresConfig = pgConfig
	}

	return s.postgresConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		grpcConfig, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = grpcConfig
	}

	return s.grpcConfig
}

func (s *serviceProvider) PostgresPoolConn(ctx context.Context) *pgxpool.Pool {
	if s.postgresPool == nil {
		pool, err := pgxpool.New(ctx, s.PostgresConfig().DSN())

		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		if err = pool.Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.postgresPool = pool
	}

	return s.postgresPool
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if s.noteRepo == nil {
		s.noteRepo = noteRepository.NewPostgresRepository(s.PostgresPoolConn(ctx))
	}
	return s.noteRepo
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(s.NoteRepository(ctx))
	}
	return s.noteService
}

func (s *serviceProvider) NoteServerImplementation(ctx context.Context) *note.ServerImplementation {
	if s.noteImpl == nil {
		s.noteImpl = note.NewServerImplementation(s.NoteService(ctx))
	}
	return s.noteImpl
}
