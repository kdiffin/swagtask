package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	db "swagtask/db/generated"
	"swagtask/handlers"
	"swagtask/middleware"
	"swagtask/models"

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
	
	
	log.SetFlags(log.LstdFlags)
	templates := models.NewTemplate()
	mux := http.NewServeMux()
	// e.Renderer = newTemplate() implement
	// .Use(middleware.Logger()) implement
	// .Static("/images", "../images")
	// e.Static("/css", "../css")
	// e.Static("/js", "../js")
	// router.Tasks(mux, queries, templates)

	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("../images"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../js"))))

	mux.HandleFunc("GET /tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetTasks(w, r, queries, templates)
	})
	mux.HandleFunc("POST /tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerCreateTask(w, r, queries, templates)
	})
	mux.HandleFunc("POST /tasks/{id}/toggle-complete/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, _ := strconv.Atoi(idStr) // f the error
		handlers.HandlerTaskToggleComplete(w, r, queries, templates, int32(id))
	})
	mux.HandleFunc("POST /tasks/{id}/tags/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, _ := strconv.Atoi(idStr) // f the error
		tagIdStr := r.FormValue("tag")
		tagId, _ := strconv.Atoi(tagIdStr) // f the error

		handlers.HandlerAddTagToTask(w,r,queries, templates, int32(id), int32(tagId))
	})
	mux.HandleFunc("DELETE /tasks/{id}/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, _ := strconv.Atoi(idStr) // f the error

		handlers.HandlerDeleteTask(w,r,queries, templates, int32(id))
	})

	

	server := http.Server{
		Addr: ":42069",
		Handler: middleware.Logging(mux),
	}
	fmt.Println("running server")
	log.Fatal(server.ListenAndServe())
	fmt.Println("not running server")

}
