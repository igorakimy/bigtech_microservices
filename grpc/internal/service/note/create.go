package note

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	id, err := s.noteRepo.Create(ctx, info)

	if err != nil {
		return 0, err
	}

	return id, nil
}
