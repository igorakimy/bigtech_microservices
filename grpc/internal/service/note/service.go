package note

import (
	"github.com/igorakimy/bigtech_microservices/internal/client/db"
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	"github.com/igorakimy/bigtech_microservices/internal/service"
)

type serv struct {
	noteRepo  repository.NoteRepository
	txManager db.TxManager
}

func NewService(
	noteRepo repository.NoteRepository,
	txManager db.TxManager,
) service.NoteService {
	return &serv{
		noteRepo:  noteRepo,
		txManager: txManager,
	}
}

func NewMockService(noteRepo repository.NoteRepository) service.NoteService {
	return &serv{
		noteRepo: noteRepo,
	}
}
