package repository

import (
	"context"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
)

type NoteRepository interface {
	Get(ctx context.Context, id int64) (*desc.Note, error)
	List(ctx context.Context) ([]*desc.Note, error)
	Create(ctx context.Context, info *desc.NoteInfo) (int64, error)
	Update(ctx context.Context, id int64, info *desc.UpdateNoteInfo) error
	Delete(ctx context.Context, id int64) error
}
