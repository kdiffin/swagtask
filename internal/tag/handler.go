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
	page := newTagsPage(tagsWithTasks, true, user.PathToPfp.String, user.Username)
	templates.Render(w, "tags-page", page)
}

func HandlerUpdateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	tagId := r.PathValue("id")
	tagName := r.FormValue("tag_name")

	if tagName == "" {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(nil))
		return
	}
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tagWithTask, err := updateTag(queries, utils.PgUUID(tagId), user.ID, tagName, r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTask)
}

// TODO:
// breaking the abstraction rules here cuz its gonna be mad annoying working with interface{}
// func HandlerCreateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template, tagName string, source string) {
// 	user, ok := middleware.UserFromContext(r.Context())
// 	if !ok {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
// 		filters := models.NewTasksPageFilters(r.URL.Query().Get("tags"), r.URL.Query().Get("taskName"))
// 		tasksWithTags, errTasks := GetFilteredTasksWithTags(queries, &filters, user.ID, r.Context())
// 		if utils.CheckError(w, r, errTasks) {
// 			return
// 		}

// 		templates.Render(w, "tasks-container", tasksWithTags)
// 		return
// 	case "/tags":
// 		// tagsWithTasks
// 		tagsWithTasks, errTags := getTagsWithTasks(queries, user.ID, r.Context())
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

func HandlerDeleteTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	tagId := r.PathValue("id")
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := deleteTag(queries, utils.PgUUID(tagId), user.ID, r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	// bc we want htmx to rerender it
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nil))
}
