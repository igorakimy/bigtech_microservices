package note

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/converter"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
)

func (s *ServerImplementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	noteID, err := s.noteService.Create(
		ctx,
		converter.ToNoteInfoFromDesc(req.GetInfo()),
	)
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: noteID}, nil
}
