package router

import (
	"net/http"
	"os"
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	"swagtask/internal/task"
	"swagtask/internal/template"
	collaborative_vault "swagtask/internal/vault/collaborative-page"
	owner_dashboard "swagtask/internal/vault/owner-dashboard"
)

func NewMux(queries *db.Queries, templates *template.Template) *http.ServeMux {
	mux := http.NewServeMux()
	taskHandler := task.NewTaskHandler(queries, templates)
	tagHandler := tag.NewTagHandler(queries, templates)
	authHandler := auth.NewAuthHandler(queries, templates)
	ownerVaultHandler := owner_dashboard.NewVaultHandler(queries, templates)
	collaborativeVaultHandler := collaborative_vault.NewVaultHandler(queries, templates)

	uploadsFS := http.FS(os.DirFS("./web/pfps/"))
	staticFS := http.FS(os.DirFS("./web/static/"))

	mux.Handle("/{$}", middleware.HandlerWithUserNoRedirect(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.UserFromContext(r.Context())

		users, _ := queries.GetUsers(r.Context())

		type pageType struct {
			Auth  auth.AuthenticatedPage
			Users []db.GetUsersRow
		}

		var page pageType
		if err != nil {
			page = pageType{
				Auth: auth.AuthenticatedPage{
					Authorized: false,
					User: auth.UserUI{
						PathToPfp: "",
						Username:  "",
					},
				},
				Users: users,
			}
		} else {
			page = pageType{
				Auth: auth.AuthenticatedPage{
					Authorized: true,
					User: auth.UserUI{
						PathToPfp: user.PathToPfp,
						Username:  user.Username,
					},
				},
				Users: users,
			}
		}
		templates.Render(w, "landing-page", page)
	})))
	mux.Handle("/pfps/", http.StripPrefix("/pfps/", http.FileServer(uploadsFS)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))

	SetupAuthRoutes(mux, authHandler)
	SetupTaskRoutes(mux, queries, taskHandler)
	SetupTagRoutes(mux, queries, tagHandler)
	SetupVaultRoutes(mux, queries, ownerVaultHandler, collaborativeVaultHandler)

	return mux
}
