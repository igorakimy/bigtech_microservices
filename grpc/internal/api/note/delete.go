package note

import (
	"context"

	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerImplementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.noteService.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
