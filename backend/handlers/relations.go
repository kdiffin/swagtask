package handlers

import (
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
)


func HandlerAddTagToTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	taskWithTags, err := service.AddTagToTask(queries, tagId, taskId)
	if checkErrors(w,err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerRemoveTagFromTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	taskWithTags, err := 	service.DeleteTagRelationFromTask(queries, tagId, taskId)
	if checkErrors(w,err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerAddTaskToTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	tagWithTasks, err := service.AddTaskToTag(queries, tagId, taskId)
	if checkErrors(w,err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTasks)
}

func HandlerRemoveTaskFromTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	tagWithTasks, err := service.DeleteTaskRelationFromTag(queries, tagId, taskId)
	if checkErrors(w, err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTasks)
}