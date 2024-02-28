package main

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const (
	dbDSN = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Create pool connection to database
	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create query for inserting row
	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.City(), gofakeit.Address().Street).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var noteID int
	if err = pool.QueryRow(ctx, query, args...).Scan(&noteID); err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	log.Printf("inserted note with id: %d", noteID)

	// Create query for select rows
	builderSelect := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var n = struct {
		id          int
		title, body string
		createdAt   time.Time
		updatedAt   sql.NullTime
	}{}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select rows: %v", err)
	}

	for rows.Next() {
		if err = rows.Scan(&n.id, &n.title, &n.body, &n.createdAt, &n.updatedAt); err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}

		log.Printf(
			"id: %d\ntitle: %s\nbody: %s\ncreated_at: %v\nupdated_at: %v\n",
			n.id,
			n.title,
			n.body,
			n.createdAt,
			n.updatedAt,
		)
	}

	// Create query for update row
	builderUpdate := sq.Update("note").
		PlaceholderFormat(sq.Dollar).
		Set("title", gofakeit.City()).
		Set("body", gofakeit.Address().Street).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": noteID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build update query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())
}
