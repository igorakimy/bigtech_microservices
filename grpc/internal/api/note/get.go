package note

import (
	"context"
	"errors"
	"github.com/igorakimy/bigtech_microservices/internal/converter"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
)

func (s *ServerImplementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
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
