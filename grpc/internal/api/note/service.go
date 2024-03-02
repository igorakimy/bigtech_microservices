package note

import (
	"github.com/igorakimy/bigtech_microservices/internal/service"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
)

type ServerImplementation struct {
	desc.UnimplementedNoteV1Server
	noteService service.NoteService
}

func NewServerImplementation(noteService service.NoteService) *ServerImplementation {
	return &ServerImplementation{
		noteService: noteService,
	}
}
