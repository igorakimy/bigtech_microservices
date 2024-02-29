package note

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/igorakimy/bigtech_microservices/internal/repository"
	"github.com/igorakimy/bigtech_microservices/internal/repository/note/converter"
	"github.com/igorakimy/bigtech_microservices/internal/repository/note/model"
	desc "github.com/igorakimy/bigtech_microservices/pkg/note/v1"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) repository.NoteRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (*desc.Note, error) {
	builder := sq.Select(idCol, titleCol, contentCol, createdAtCol, updatedAtCol).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idCol: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var n model.Note
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&n.ID, &n.Info.Title, &n.Info.Body, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&n), nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]*desc.Note, error) {
	builder := sq.Select(idCol, titleCol, contentCol, createdAtCol, updatedAtCol).
		PlaceholderFormat(sq.Dollar).
		From(tableName)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var n model.Note
	var notes []model.Note

	for rows.Next() {
		err = rows.Scan(&n.ID, &n.Info.Title, &n.Info.Body, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return converter.ToNotesFromRepo(notes), nil
}

func (r *PostgresRepository) Create(ctx context.Context, info *desc.NoteInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleCol, contentCol, createdAtCol).
		Values(info.GetTitle(), info.GetContent(), time.Now()).
		Suffix("RETURNING " + idCol)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var noteID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&noteID)
	if err != nil {
		return 0, err
	}

	return noteID, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id int64, info *desc.UpdateNoteInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(titleCol, info.GetTitle().GetValue()).
		Set(contentCol, info.GetContent().GetValue()).
		Set(updatedAtCol, time.Now()).
		Where(sq.Eq{idCol: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
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

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
