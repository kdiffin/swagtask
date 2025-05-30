package router

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/task"
	template "swagtask/internal/template"
)

func SetupTaskRoutes(mux *http.ServeMux, queries *db.Queries, template *template.Template) {

	mux.Handle("GET /tasks/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task.HandlerGetTasks(w, r, queries, template)
	})))
	mux.Handle("GET /tasks/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		task.HandlerGetTask(w, r, queries, template)
	})))
	mux.Handle("POST /tasks/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task.HandlerCreateTask(w, r, queries, template)
	})))
	mux.Handle("POST /tasks/{id}/toggle-complete/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task.HandlerTaskToggleComplete(w, r, queries, template)
	})))
	mux.Handle("POST /tasks/{id}/tags/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task.HandlerAddTagToTask(w, r, queries, template)
	})))
	mux.Handle("DELETE /tasks/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task.HandlerDeleteTask(w, r, queries, template)
	})))
	mux.Handle("DELETE /tasks/{id}/tags/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		task.HandlerRemoveTagFromTask(w, r, queries, template)
	})))
	mux.Handle("PUT /tasks/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task.HandlerUpdateTask(w, r, queries, template)
	})))

}
