package note

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/igorakimy/bigtech_microservices/internal/client/db"
	"github.com/igorakimy/bigtech_microservices/internal/model"
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	"github.com/igorakimy/bigtech_microservices/internal/repository/note/converter"
	modelRepo "github.com/igorakimy/bigtech_microservices/internal/repository/note/model"
	"time"
)

const (
	tableName = "note"

	idCol        = "id"
	titleCol     = "title"
	contentCol   = "body"
	createdAtCol = "created_at"
	updatedAtCol = "updated_at"
)

type PostgresRepository struct {
	db db.Client
}

func NewPostgresRepository(db db.Client) repository.NoteRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (*model.Note, error) {
	builder := sq.Select(idCol, titleCol, contentCol, createdAtCol, updatedAtCol).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idCol: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "note_repository.Get",
		QueryRaw: query,
	}

	var n modelRepo.Note
	err = r.db.DB().ScanOneContext(ctx, &n, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&n), nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]model.Note, error) {
	builder := sq.Select(idCol, titleCol, contentCol, createdAtCol, updatedAtCol).
		PlaceholderFormat(sq.Dollar).
		From(tableName)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "note_repository.List",
		QueryRaw: query,
	}

	var notes []modelRepo.Note

	err = r.db.DB().ScanAllContext(ctx, &notes, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToNotesFromRepo(notes), nil
}

func (r *PostgresRepository) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleCol, contentCol, createdAtCol).
		Values(info.Title, info.Body, time.Now()).
		Suffix("RETURNING " + idCol)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "note_repository.Create",
		QueryRaw: query,
	}

	var noteID int64
	err = r.db.DB().ScanOneContext(ctx, &noteID, q, args...)
	if err != nil {
		return 0, err
	}

	return noteID, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id int64, info *model.UpdateNoteInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(titleCol, info.Title).
		Set(contentCol, info.Body).
		Set(updatedAtCol, time.Now()).
		Where(sq.Eq{idCol: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "note_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idCol: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "note_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
