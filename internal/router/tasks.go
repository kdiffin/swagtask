package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/task"
)

func SetupTaskRoutes(mux *http.ServeMux, queries *db.Queries, handler *task.TaskHandler) {

	mux.Handle("GET /tasks/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.GetAll)))
	mux.Handle("GET /tasks/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.GetByID)))
	mux.Handle("POST /tasks/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.Create)))
	mux.Handle("POST /tasks/{id}/toggle-complete/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.ToggleComplete)))
	mux.Handle("POST /tasks/{id}/tags/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.AddTag)))
	mux.Handle("DELETE /tasks/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.Delete)))
	mux.Handle("DELETE /tasks/{id}/tags/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.RemoveTag)))
	mux.Handle("PUT /tasks/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(handler.Update)))

}
