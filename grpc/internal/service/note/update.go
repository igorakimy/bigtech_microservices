package note

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/model"
)

func (s *serv) Update(ctx context.Context, id int64, info *model.UpdateNoteInfo) error {
	if err := s.noteRepo.Update(ctx, id, info); err != nil {
		return err
	}
	return nil
}
