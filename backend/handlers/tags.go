package handlers

import (
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
)
 
func HandlerGetTags(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template) {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r,errAuth)  {
		return
	}
	tagsWithTasks, errTag := service.GetTagsWithTasks(queries, user.ID, r.Context())
	if checkErrors(w,r,errTag) {
		return
	}
	page := models.NewTagsPage(tagsWithTasks, true, user.PathToPfp, user.Username)
	templates.Render(w, "tags-page", page)
}


func HandlerUpdateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, tagName string, tagId int32) {
	if tagName == "" {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(nil))
		return
	}
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	tagWithTask, err := service.UpdateTag(queries, tagId, user.ID, tagName, r.Context())
	if checkErrors(w,r, err) {
		return 
	}

	templates.Render(w, "tag-card", tagWithTask)
}
	

// breaking the abstraction rules here cuz its gonna be mad annoying working with interface{}
func HandlerCreateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *models.Template, tagName string, source string)  {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	err := queries.CreateTag(r.Context(), db.CreateTagParams{
		Name: tagName,
		UserID: user.ID,
	})
	if checkErrors(w,r,err) {
		return
	}

	switch source {
	case "/tasks":
		filters := models.NewTasksPageFilters(r.URL.Query().Get("tags"), r.URL.Query().Get("taskName"))
		tasksWithTags, errTasks := service.GetFilteredTasksWithTags(queries, &filters, user.ID, r.Context())
		if checkErrors(w,r,errTasks) {
			return
		}
		
		templates.Render(w, "tasks-container", tasksWithTags)
		return 
	case "/tags":
		// tagsWithTasks
		tagsWithTasks, errTags := service.GetTagsWithTasks(queries, user.ID, r.Context())
		if checkErrors(w,r,errTags) {
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
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	err := service.DeleteTag(queries, tagId,user.ID, r.Context())
	if checkErrors(w,r,err) {
		return
	}

	// bc we want htmx to rerender it
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nil))
}