package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	collaborative_vault "swagtask/internal/vault/collaborative-page"
	owner_dashboard "swagtask/internal/vault/owner-dashboard"

	"golang.org/x/net/websocket"
)

func SetupVaultRoutes(
	mux *http.ServeMux,
	queries *db.Queries,
	ownerHandler *owner_dashboard.VaultHandler,
	collaborativeHandler *collaborative_vault.VaultHandler,
) {
	// vault page
	mux.Handle("GET /vaults/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(ownerHandler.GetAll)))

	mux.Handle("POST /vaults/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(ownerHandler.Create)))
	mux.Handle("DELETE /vaults/{vaultId}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(ownerHandler.Delete)))
	mux.Handle("PUT /vaults/{vaultId}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(ownerHandler.Update)))
	mux.Handle("POST /vaults/{vaultId}/collaborators/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(ownerHandler.AddCollaborator)))
	mux.Handle("DELETE /vaults/{vaultId}/collaborators/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(ownerHandler.RemoveCollaborator)))

	// all dynamic pages should have vaultId as their parameter
	// ^ and all pages should be protected
	// ^ check how the middleware works
	// vaults{id} page with collaboration
	mux.Handle("GET /vaults/{vaultId}/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(collaborativeHandler.Get))))
	mux.Handle("GET /vaults/{vaultId}/tasks/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(collaborativeHandler.GetTasks))))
	mux.Handle("GET /vaults/{vaultId}/tasks/{id}/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(collaborativeHandler.GetTask))))
	mux.Handle("GET /vaults/{vaultId}/tags/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(collaborativeHandler.GetTags))))

	mux.Handle("/vaults/{vaultId}/wstest/{$}", middleware.HandlerWithVaultIdFromPath(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vaultId, _ := middleware.VaultIDFromContext(r.Context())
			websocket.Handler(collaborative_vault.WsHandlerPubSubTest(vaultId)).ServeHTTP(w, r)
		},
	)))

	mux.Handle("/vaults/{vaultId}/ws/{$}", middleware.HandlerWithVaultIdFromPath(middleware.HandlerWithUser(queries, http.HandlerFunc(collaborativeHandler.WebSocket))))

	mux.HandleFunc("/debug", collaborative_vault.DebugHandler()) // <<== here

}
