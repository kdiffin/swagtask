package router

import (
	"net/http"
	"swagtask/backend/handlers"
	db "swagtask/internal/db/generated"
	"swagtask/internal/template"

	"golang.org/x/net/websocket"
)

func SetupVaultRoutes(mux *http.ServeMux, queries *db.Queries, templates *template.Template) {
	mux.HandleFunc("GET /vaults/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetVaults(w, r, queries, templates)
	})
	mux.Handle("/ws/", websocket.Handler(handlers.WsHandler()))
}
