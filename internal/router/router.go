package router

import (
	"net/http"
	"os"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	"swagtask/internal/template"
)

func NewMux(queries *db.Queries, templates *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	uploadsFS := http.FS(os.DirFS("./web/pfps/"))
	staticFS := http.FS(os.DirFS("./web/static/"))

	mux.Handle("/pfps/", http.StripPrefix("/pfps/", http.FileServer(uploadsFS)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hellow orld"))
	})
	mux.Handle("POST /tags/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerCreateTag(w, r, queries, templates)
	})))
	SetupAuthRoutes(mux, queries, templates)
	SetupTaskRoutes(mux, queries, templates)
	SetupTagRoutes(mux, queries, templates)
	SetupVaultRoutes(mux, queries, templates)

	return mux
}
