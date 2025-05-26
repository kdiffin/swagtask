package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	collaborative_vault "swagtask/internal/vault/collaborative-page"
	owner_dashboard "swagtask/internal/vault/owner-dashboard"
)

func SetupVaultRoutes(mux *http.ServeMux, queries *db.Queries, templates *template.Template) {
	// vault page
	mux.Handle("GET /vaults/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		owner_dashboard.HandlerGetVaults(w, r, queries, templates)
	})))

	mux.Handle("POST /vaults/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		owner_dashboard.HandlerCreateVault(w, r, queries, templates)
	})))
	mux.Handle("DELETE /vaults/{vaultId}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		owner_dashboard.HandlerDeleteVault(w, r, queries, templates)
	})))
	mux.Handle("PUT /vaults/{vaultId}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		owner_dashboard.HandlerUpdateVault(w, r, queries, templates)
	})))
	mux.Handle("POST /vaults/{vaultId}/collaborators/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		owner_dashboard.HandlerAddCollaboratorToVault(w, r, queries, templates)
	})))
	mux.Handle("DELETE /vaults/{vaultId}/collaborators/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		owner_dashboard.HandlerDeleteCollaboratorToVault(w, r, queries, templates)
	})))

	// all dynamic pages should have vaultId as their parameter
	// ^ and all pages should be protected
	// ^ check how the middleware works
	// vaults{id} page with collaboration
	mux.Handle("GET /vaults/{vaultId}/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		collaborative_vault.HandlerGetVault(w, r, queries, templates)
	}))))
	mux.Handle("GET /vaults/{vaultId}/tasks/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		collaborative_vault.HandlerGetTasks(w, r, queries, templates)
	}))))
	// mux.Handle("GET /vaults/{vaultId}/tags/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// collaborative_vault.HandlerGetTags(w, r, queries, templates)
	// }))))
	// mux.Handle("/ws/", websocket.Handler(handlers.WsHandler()))
}
