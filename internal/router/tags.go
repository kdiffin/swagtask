package router

import (
	"encoding/json"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

func SetupTagRoutes(mux *http.ServeMux, queries *db.Queries, templates *template.Template) {
	// mux.HandleFunc("POST /tags/{$}", func(w http.ResponseWriter, r *http.Request) {
	// 	tag.HandlerCreateTag(w, r, queries, templates)
	// })
	mux.Handle("GET /tags/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerGetTags(w, r, queries, templates)
	})))
	mux.Handle("PUT /tags/{id}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerUpdateTag(w, r, queries, templates)
	})))
	mux.Handle("DELETE /tags/{id}/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerDeleteTag(w, r, queries, templates)
	})))
	// mux.HandleFunc("POST /tags/{id}/tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := r.PathValue("id")
	// 	id, errConv := strconv.Atoi(idStr)
	// 	taskIdStr := r.FormValue("task_id")
	// 	taskId, errConv2 := strconv.Atoi(taskIdStr)
	// 	if errConv != nil || errConv2 != nil {
	// 		utils.LogError("couldnt convert to str", errConv)
	// 		utils.LogError("couldnt convert to str2", errConv2)
	// 		return
	// 	}

	// 	tag.(w, r, queries, templates, int32(taskId), int32(id))
	// })
	// mux.HandleFunc("DELETE /tags/{id}/tasks/{$}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := r.PathValue("id")
	// 	id, errConv := strconv.Atoi(idStr)
	// 	taskIdStr := r.FormValue("task_id")
	// 	taskId, errConv2 := strconv.Atoi(taskIdStr)
	// 	if errConv != nil || errConv2 != nil {
	// 		utils.LogError("couldnt convert to str", errConv)
	// 		utils.LogError("couldnt convert to str2", errConv2)
	// 		return
	// 	}

	// 	tag.HandlerRemoveTaskFromTag(w, r, queries, templates, int32(taskId), int32(id))
	// })
	mux.HandleFunc("GET /json", func(w http.ResponseWriter, r *http.Request) {
		// Prepare the response data
		response := map[string]string{
			"message": "hello world",
		}

		// Set content-type to JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode and send JSON
		if err := json.NewEncoder(w).Encode(response); err != nil {
			utils.LogError("failed to encode json", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})

}
