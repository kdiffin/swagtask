package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/template"
)

func SetupVaultRoutes(mux *http.ServeMux, queries *db.Queries, templates *template.Template) {
	// mux.HandleFunc("GET /vaults/{$}", func(w http.ResponseWriter, r *http.Request) {
	// 	vault.HandlerGetVaults(w, r, queries, templates)
	// })
	// mux.Handle("/ws/", websocket.Handler(handlers.WsHandler()))
}
