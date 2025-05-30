package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/router"
	"swagtask/internal/template"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// DB INIT START
	// pgx pool starts a pool thats concurrency safe
	dbpool, errDbConn := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	// connects to db via url
	if errDbConn != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", errDbConn)
		os.Exit(1)
	}
	defer dbpool.Close()
	queries := db.New(dbpool)
	// DB INIT END

	templates := template.NewTemplate()
	mux := router.NewMux(queries, templates)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}
	server := http.Server{
		Addr:    ":" + port,
		Handler: middleware.Logging(mux),
	}

	fmt.Println("running server")

	log.Fatal(server.ListenAndServe())
	fmt.Println("not running servers")

}
