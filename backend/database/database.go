package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DatabaseInit() *pgxpool.Pool {
	// pgx pool starts a pool thats concurrency safe
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	// connects to db via url
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
