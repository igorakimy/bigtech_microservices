package repository

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/model"
)

type NoteRepository interface {
	Get(ctx context.Context, id int64) (*model.Note, error)
	List(ctx context.Context) ([]model.Note, error)
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Update(ctx context.Context, id int64, info *model.UpdateNoteInfo) error
	Delete(ctx context.Context, id int64) error
}
