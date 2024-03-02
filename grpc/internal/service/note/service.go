package note

import (
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	"github.com/igorakimy/bigtech_microservices/internal/service"
)

type serv struct {
	noteRepo repository.NoteRepository
}

func NewService(noteRepo repository.NoteRepository) service.NoteService {
	return &serv{
		noteRepo: noteRepo,
	}
}
