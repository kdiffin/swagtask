package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/task"
)

func SetupTaskRoutes(mux *http.ServeMux, queries *db.Queries, handler *task.TaskHandler) {
	withVault := func(next http.HandlerFunc) http.Handler {
		return middleware.HandlerWithVaultIdFromUser(queries, next)
	}

	mux.Handle("GET /tasks/{$}", withVault(handler.GetAll))
	mux.Handle("GET /tasks/{id}/{$}", withVault(handler.GetByID))
	mux.Handle("POST /tasks/{$}", withVault(handler.Create))
	mux.Handle("POST /tasks/tag-options/{$}", withVault(handler.CreateTagOption))
	mux.Handle("POST /tasks/{id}/toggle-complete/{$}", withVault(handler.ToggleComplete))
	mux.Handle("POST /tasks/{id}/tags/{$}", withVault(handler.AddTag))
	mux.Handle("DELETE /tasks/{id}/{$}", withVault(handler.Delete))
	mux.Handle("DELETE /tasks/{id}/tags/{$}", withVault(handler.RemoveTag))
	mux.Handle("PUT /tasks/{id}/{$}", withVault(handler.Update))

}
