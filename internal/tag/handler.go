package tag

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

func HandlerGetTags(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tagsWithTasks, errTag := getTagsWithTasks(queries, user.ID, r.Context())
	if utils.CheckError(w, r, errTag) {
		return
	}
	page := newTagsPage(tagsWithTasks, true, user.PathToPfp, user.Username)
	templates.Render(w, "tags-page", page)
}

// func HandlerUpdateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template, tagName string, tagId int32) {
// 	if tagName == "" {
// 		w.WriteHeader(http.StatusNoContent)
// 		w.Write([]byte(nil))
// 		return
// 	}
// 	user, errAuth := getUserInfoFromSessionId(queries, r)
// 	if utils.CheckError(w, r, errAuth) {
// 		return
// 	}
// 	tagWithTask, err := updateTag(queries, tagId, user.ID, tagName, r.Context())
// 	if utils.CheckError(w, r, err) {
// 		return
// 	}

// 	templates.Render(w, "tag-card", tagWithTask)
// }

// breaking the abstraction rules here cuz its gonna be mad annoying working with interface{}
// func HandlerCreateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template, tagName string, source string) {
// 	user, errAuth := getUserInfoFromSessionId(queries, r)
// 	if utils.CheckError(w, r, errAuth) {
// 		return
// 	}
// 	err := queries.CreateTag(r.Context(), db.CreateTagParams{
// 		Name:   tagName,
// 		UserID: user.ID,
// 	})
// 	if utils.CheckError(w, r, err) {
// 		return
// 	}

// 	switch source {
// 	case "/tasks":
// 		filters := template.NewTasksPageFilters(r.URL.Query().Get("tags"), r.URL.Query().Get("taskName"))
// 		tasksWithTags, errTasks := GetFilteredTasksWithTags(queries, &filters, user.ID, r.Context())
// 		if utils.CheckError(w, r, errTasks) {
// 			return
// 		}

// 		templates.Render(w, "tasks-container", tasksWithTags)
// 		return
// 	case "/tags":
// 		// tagsWithTasks
// 		tagsWithTasks, errTags := GetTagsWithTasks(queries, user.ID, r.Context())
// 		if utils.CheckError(w, r, errTags) {
// 			return
// 		}

// 		templates.Render(w, "tags-list-container", tagsWithTasks)
// 		return
// 	default:
// 		http.Error(w, "what u on bruh", http.StatusBadGateway)
// 		return
// 	}
// }

// func HandlerDeleteTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template, tagId int32) {
// 	user, errAuth := getUserInfoFromSessionId(queries, r)
// 	if utils.CheckError(w, r, errAuth) {
// 		return
// 	}
// 	err := DeleteTag(queries, tagId, user.ID, r.Context())
// 	if utils.CheckError(w, r, err) {
// 		return
// 	}

// 	// bc we want htmx to rerender it
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(nil))
// }
