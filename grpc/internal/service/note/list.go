package note

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/model"
)

func (s *serv) List(ctx context.Context) ([]model.Note, error) {
	_, err := s.noteRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
