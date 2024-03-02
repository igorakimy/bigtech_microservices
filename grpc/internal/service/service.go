package service

import (
	"context"
	"github.com/igorakimy/bigtech_microservices/internal/model"
)

type NoteService interface {
	Create(ctx context.Context, info *model.NoteInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Note, error)
	List(ctx context.Context) ([]model.Note, error)
	Update(ctx context.Context, id int64, info *model.UpdateNoteInfo) error
	Delete(ctx context.Context, id int64) error
}
