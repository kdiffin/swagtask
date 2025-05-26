package router

import (
	"fmt"
	"net/http"
	"os"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

func NewMux(queries *db.Queries, templates *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	uploadsFS := http.FS(os.DirFS("./web/pfps/"))
	staticFS := http.FS(os.DirFS("./web/static/"))

	mux.Handle("/pfps/", http.StripPrefix("/pfps/", http.FileServer(uploadsFS)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hellow orld"))
	})
	mux.Handle("POST /tags/{$}", middleware.HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandlerCreateTag(w, r, queries, templates)
	})))
	SetupAuthRoutes(mux, queries, templates)
	SetupTaskRoutes(mux, queries, templates)
	SetupTagRoutes(mux, queries, templates)
	SetupVaultRoutes(mux, queries, templates)

	return mux
}

// / breaking the abstraction rules here cuz its gonna be mad annoying working with interface{}
func HandlerCreateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	tagName := r.FormValue("tag_name")
	source := r.FormValue("source")
	user, errUser := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, errUser) {
		return
	}

	fmt.Println(user.DefaultVaultID)
	err := queries.CreateTag(r.Context(), db.CreateTagParams{
		Name:    tagName,
		UserID:  utils.PgUUID(user.ID),
		VaultID: utils.PgUUID(user.DefaultVaultID),
	})
	if utils.CheckError(w, r, err) {
		fmt.Println("error was here1")
		return
	}

	switch source {
	// case "/tasks":
	// 	filters := models.NewTasksPageFilters(r.URL.Query().Get("tags"), r.URL.Query().Get("taskName"))
	// 	tasksWithTags, errTasks := task.GetFilteredTasksWithTags(queries, &filters, user.ID, r.Context())
	// 	if utils.CheckError(w, r, errTasks) {
	// 		return
	// 	}

	// 	templates.Render(w, "tasks-container", tasksWithTags)
	// 	return
	case "/tags":
		// tagsWithTasks
		tagsWithTasks, errTags := tag.GetTagsWithTasks(queries, utils.PgUUID(user.ID), utils.PgUUID(user.DefaultVaultID), r.Context())
		if utils.CheckError(w, r, errTags) {
			fmt.Println("error was here")
			return
		}

		templates.Render(w, "tags-list-container", tagsWithTasks)
		return
	default:
		http.Error(w, "what u on bruh", http.StatusBadGateway)
		return
	}
}
