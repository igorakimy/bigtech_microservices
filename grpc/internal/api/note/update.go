package note

import (
	"context"

	"github.com/igorakimy/bigtech_microservices/internal/converter"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerImplementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
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
