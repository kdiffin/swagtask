package router

import (
	"net/http"
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	template "swagtask/internal/template"
)

func SetupAuthRoutes(mux *http.ServeMux, queries *db.Queries, template *template.Template) {
	mux.HandleFunc("POST /sign-up/{$}", func(w http.ResponseWriter, r *http.Request) {
		auth.HandleSignup(queries, w, r)
	})
	mux.HandleFunc("GET /sign-up/{$}", func(w http.ResponseWriter, r *http.Request) {
		template.Render(w, "sign-up", nil)
	})
	mux.HandleFunc("POST /login/{$}", func(w http.ResponseWriter, r *http.Request) {
		auth.HandleLogin(queries, w, r)
	})
	mux.HandleFunc("GET /login/{$}", func(w http.ResponseWriter, r *http.Request) {
		template.Render(w, "login", nil)
	})
	mux.HandleFunc("POST /logout/{$}", func(w http.ResponseWriter, r *http.Request) {
		auth.HandleLogout(queries, w, r)
	})
}
