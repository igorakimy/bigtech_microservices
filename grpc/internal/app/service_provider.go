package app

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/api/note"
	"github.com/igorakimy/bigtech_microservices/internal/client/db"
	"github.com/igorakimy/bigtech_microservices/internal/client/db/pg"
	"github.com/igorakimy/bigtech_microservices/internal/client/db/transaction"
	"github.com/igorakimy/bigtech_microservices/internal/closer"
	"github.com/igorakimy/bigtech_microservices/internal/config"
	"github.com/igorakimy/bigtech_microservices/internal/config/env"
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	noteRepository "github.com/igorakimy/bigtech_microservices/internal/repository/note"
	"github.com/igorakimy/bigtech_microservices/internal/service"
	noteService "github.com/igorakimy/bigtech_microservices/internal/service/note"
	"log"
)

type serviceProvider struct {
	postgresConfig config.PostgresConfig
	grpcConfig     config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager
	noteRepo  repository.NoteRepository

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

func (s *serviceProvider) DbClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		dbClient, err := pg.New(ctx, s.PostgresConfig().DSN())

		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		if err = dbClient.DB().Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(dbClient.Close)

		s.dbClient = dbClient
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DbClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if s.noteRepo == nil {
		s.noteRepo = noteRepository.NewPostgresRepository(s.DbClient(ctx))
	}
	return s.noteRepo
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(
			s.NoteRepository(ctx),
			s.TxManager(ctx),
		)
	}
	return s.noteService
}

func (s *serviceProvider) NoteServerImplementation(ctx context.Context) *note.ServerImplementation {
	if s.noteImpl == nil {
		s.noteImpl = note.NewServerImplementation(s.NoteService(ctx))
	}
	return s.noteImpl
}
