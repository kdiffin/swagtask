package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/vault"
)

func SetupVaultRoutes(mux *http.ServeMux, queries *db.Queries, templates *template.Template) {
	mux.Handle("GET /vaults/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vault.HandlerGetVaults(w, r, queries, templates)
	})))
	mux.Handle("POST /vaults/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vault.HandlerCreateVault(w, r, queries, templates)
	})))
	mux.Handle("DELETE /vaults/{vaultId}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vault.HandlerDeleteVault(w, r, queries, templates)
	})))
	mux.Handle("PUT /vaults/{vaultId}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vault.HandlerUpdateVault(w, r, queries, templates)
	})))
	mux.Handle("POST /vaults/{vaultId}/collaborators/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vault.HandlerAddCollaboratorToVault(w, r, queries, templates)
	})))
	mux.Handle("DELETE /vaults/{vaultId}/collaborators/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vault.HandlerDeleteCollaboratorToVault(w, r, queries, templates)
	})))
	// mux.Handle("/ws/", websocket.Handler(handlers.WsHandler()))
}
