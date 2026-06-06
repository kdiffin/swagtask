package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
)

func SetupTagRoutes(mux *http.ServeMux, queries *db.Queries, handler *tag.TagHandler) {
	withVault := func(next http.HandlerFunc) http.Handler {
		return middleware.HandlerWithVaultIdFromUser(queries, next)
	}

	mux.Handle("GET /tags/{$}", withVault(handler.GetAll))
	mux.Handle("POST /tags/{$}", withVault(handler.Create))
	mux.Handle("PUT /tags/{id}/{$}", withVault(handler.Update))

	mux.Handle("DELETE /tags/{id}/{$}", withVault(handler.Delete))
	mux.Handle("POST /tags/{id}/tasks/{$}", withVault(handler.AddTask))
	mux.Handle("DELETE /tags/{id}/tasks/{$}", withVault(handler.RemoveTask))

}
