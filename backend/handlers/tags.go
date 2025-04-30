package handlers

import (
	"context"
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
)
 
func HandlerGetTags(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template) {
	tagsWithTasks, errTag := service.GetTagsWithTasks(queries)
	if checkErrors(w,errTag) {
		return
	}
	page := models.NewTagsPage(tagsWithTasks)
	templates.Render(w, "tags-page", page)
}


func HandlerUpdateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, tagName string, tagId int32) {
	if tagName == "" {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(nil))
		return
	}

	tagWithTask, err := service.UpdateTag(queries, tagId, tagName)
	if checkErrors(w, err) {
		return 
	}

	templates.Render(w, "tag-card", tagWithTask)
}
	

// breaking the abstraction rules here cuz its gonna be mad annoying working with interface{}
func HandlerCreateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, tagName string, source string)  {
	err := queries.CreateTag(context.Background(), tagName)
	if checkErrors(w,err) {
		return
	}

	switch source {
	case "/tasks":
		tasksWithTags, errTasks := service.GetTasksWithTags(queries)
		if checkErrors(w,errTasks) {
			return
		}
		
		page := models.NewTasksPage(tasksWithTags)
		templates.Render(w, "tasks-container", page)
		return 
	case "/tags":
		// tagsWithTasks
		tagsWithTasks, errTags := service.GetTagsWithTasks(queries)
		if checkErrors(w,errTags) {
			return
		}
		
		templates.Render(w, "tags-list-container", tagsWithTasks)
		return 
	default:
		http.Error(w, "what u on bruh", http.StatusBadGateway)
		return  
	}
}


func HandlerDeleteTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, tagId int32) {
	err := service.DeleteTag(queries, tagId)
	if checkErrors(w,err) {
		return
	}

	// bc we want htmx to rerender it
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nil))
}