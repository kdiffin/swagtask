package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/template"
)

func NewMux(queries *db.Queries, templates *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./web/images"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hellow orld"))
	})
	SetupAuthRoutes(mux, queries, templates)
	// SetupTaskRoutes(mux, queries, templates)
	SetupTagRoutes(mux, queries, templates)
	SetupVaultRoutes(mux, queries, templates)

	return mux
}
