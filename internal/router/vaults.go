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
	withUser := func(next http.HandlerFunc) http.Handler {
		return middleware.HandlerWithUser(queries,next)
	}
	// vault page
	mux.Handle("GET /vaults/{$}", withUser(ownerHandler.GetAll))

	mux.Handle("POST /vaults/{$}", withUser(ownerHandler.Create))
	mux.Handle("DELETE /vaults/{vaultId}/{$}", withUser(ownerHandler.Delete))
	mux.Handle("PUT /vaults/{vaultId}/{$}", withUser(ownerHandler.Update))
	mux.Handle("POST /vaults/{vaultId}/collaborators/{$}", withUser(ownerHandler.AddCollaborator))
	mux.Handle("DELETE /vaults/{vaultId}/collaborators/{$}", withUser(ownerHandler.RemoveCollaborator))

	// all dynamic pages should have vaultId as their parameter
	// ^ and all pages should be protected
	// ^ check how the middleware works
	// vaults{id} page with collaboration
	mux.Handle("GET /vaults/{vaultId}/{$}", withUser(collaborativeHandler.Get))
	mux.Handle("GET /vaults/{vaultId}/tasks/{$}", withUser(collaborativeHandler.GetTasks))
	mux.Handle("GET /vaults/{vaultId}/tasks/{id}/{$}", withUser(collaborativeHandler.GetTask))
	mux.Handle("GET /vaults/{vaultId}/tags/{$}", withUser(collaborativeHandler.GetTags))

	// REFACTOR THIS
	mux.Handle("/vaults/{vaultId}/wstest/{$}", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vaultId :=  r.PathValue("vaultId")
			websocket.Handler(collaborative_vault.WsHandlerPubSubTest(vaultId)).ServeHTTP(w, r)
		},
	))

	mux.Handle("/vaults/{vaultId}/ws/{$}", withUser(collaborativeHandler.WebSocket))

	mux.HandleFunc("/debug", collaborative_vault.DebugHandler()) // <<== here

}
