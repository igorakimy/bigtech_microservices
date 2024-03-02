package note

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/converter"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
)

func (s *ServerImplementation) List(ctx context.Context, _ *desc.ListRequest) (*desc.ListResponse, error) {
	notes, err := s.noteService.List(ctx)
	if err != nil {
		return nil, err
	}
	return &desc.ListResponse{
		Notes: converter.ToNotesFromService(notes),
	}, nil
}
