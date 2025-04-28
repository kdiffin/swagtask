package handlers

import (
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
	"swagtask/utils"
)


func HandlerAddTagToTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	taskWithTags, err := service.AddTagToTask(queries, tagId, taskId)
	if err != nil {
		utils.LogError("Updating Task completion", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerRemoveTagFromTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	taskWithTags, err := 	service.DeleteTagRelationFromTask(queries, tagId, taskId)
	if err != nil {
		utils.LogError("Updating Task completion", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	templates.Render(w, "task", taskWithTags)
}
