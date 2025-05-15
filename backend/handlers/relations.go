package handlers

import (
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
)

func HandlerAddTagToTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w, r, errAuth) {
		return
	}
	taskWithTags, err := service.AddTagToTask(queries, tagId, user.ID, taskId, r.Context())
	if checkErrors(w, r, err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerRemoveTagFromTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w, r, errAuth) {
		return
	}
	taskWithTags, err := service.DeleteTagRelationFromTask(queries, tagId, user.ID, taskId, r.Context())
	if checkErrors(w, r, err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerAddTaskToTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w, r, errAuth) {
		return
	}
	tagWithTasks, err := service.AddTaskToTag(queries, tagId, user.ID, taskId, r.Context())
	if checkErrors(w, r, err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTasks)
}

func HandlerRemoveTaskFromTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w, r, errAuth) {
		return
	}
	tagWithTasks, err := service.DeleteTaskRelationFromTag(queries, tagId, user.ID, taskId, r.Context())
	if checkErrors(w, r, err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTasks)
}
