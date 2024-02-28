package main

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

const (
	dbDSN = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Create database connection
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = con.Close(ctx) }()

	res, err := con.Exec(
		ctx,
		"INSERT INTO note (title, body) VALUES ($1, $2)",
		gofakeit.City(),
		gofakeit.Address().Street,
	)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}
	log.Printf("inserted %d rows", res.RowsAffected())

	// Select rows from database
	rows, err := con.Query(
		ctx,
		"SELECT id, title, body, created_at, updated_at FROM note",
	)
	defer rows.Close()

	for rows.Next() {
		var s = struct {
			id          int
			title, body string
			createdAt   time.Time
			updatedAt   sql.NullTime
		}{}

		err = rows.Scan(&s.id, &s.title, &s.body, &s.createdAt, &s.updatedAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}

		log.Printf(
			"id: %d\ntitle: %s\nbody: %s\ncreated:%v\nupdated:%v\n",
			s.id,
			s.title,
			s.body,
			s.createdAt,
			s.updatedAt,
		)
	}
}
