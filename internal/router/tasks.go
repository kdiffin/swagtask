package router

import (
	"net/http"
	"strconv"
	"swagtask/backend/handlers"
	db "swagtask/internal/db/generated"
	template "swagtask/internal/template"
	"swagtask/internal/utils"
)

func SetupTaskRoutes(mux *http.ServeMux, queries *db.Queries, template *template.Template) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/tasks/", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerGetTasks(w, r, queries, template)
	})
	mux.HandleFunc("GET /tasks/{id}/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr)
		if errConv != nil {
			utils.LogError("couldnt convert to str", errConv)
			return
		}

		handlers.HandlerGetTask(w, r, queries, template, int32(id))
	})
	mux.HandleFunc("POST /tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerCreateTask(w, r, queries, template)
	})
	mux.HandleFunc("POST /tasks/{id}/toggle-complete/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr)
		if errConv != nil {
			utils.LogError("couldnt convert to str", errConv)
			return
		}

		handlers.HandlerTaskToggleComplete(w, r, queries, template, int32(id))
	})
	mux.HandleFunc("POST /tasks/{id}/tags/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv1 := strconv.Atoi(idStr)
		tagIdStr := r.FormValue("tag_id")
		tagId, errConv := strconv.Atoi(tagIdStr)
		if errConv != nil || errConv1 != nil {
			utils.LogError("couldnt convert to str", errConv)
			utils.LogError("couldnt convert to str", errConv1)
			return
		}

		handlers.HandlerAddTagToTask(w, r, queries, template, int32(id), int32(tagId))
	})
	mux.HandleFunc("DELETE /tasks/{id}/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr)
		if errConv != nil {
			utils.LogError("couldnt convert to str", errConv)
			return
		}

		handlers.HandlerDeleteTask(w, r, queries, template, int32(id))
	})
	mux.HandleFunc("DELETE /tasks/{id}/tags/{$}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr)
		tagIdStr := r.FormValue("tag_id")
		tagId, errConv1 := strconv.Atoi(tagIdStr)
		if errConv != nil || errConv1 != nil {
			utils.LogError("couldnt convert to str", errConv)
			utils.LogError("couldnt convert to str 2nd", errConv1)
			return
		}

		handlers.HandlerRemoveTagFromTask(w, r, queries, template, int32(id), int32(tagId))
	})
	mux.HandleFunc("PUT /tasks/{id}/", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("task_name")
		idea := r.FormValue("task_idea")
		idStr := r.PathValue("id")
		id, errConv := strconv.Atoi(idStr)
		if errConv != nil {
			utils.LogError("couldnt convert to str", errConv)
			return
		}
		handlers.HandlerUpdateTask(w, r, queries, template, int32(id), idea, name)
	})

}
