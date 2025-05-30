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
)

func NewMux(queries *db.Queries, templates *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	uploadsFS := http.FS(os.DirFS("./web/pfps/"))
	staticFS := http.FS(os.DirFS("./web/static/"))

	mux.Handle("/{$}", middleware.HandlerWithUserNoRedirect(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.UserFromContext(r.Context())

		type pageType struct {
			Auth auth.AuthenticatedPage
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
			}
		}
		templates.Render(w, "landing-page", page)
	})))
	mux.Handle("/pfps/", http.StripPrefix("/pfps/", http.FileServer(uploadsFS)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))

	mux.Handle("POST /tags/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("source") == "/tags" {
			tag.HandlerCreateTag(w, r, queries, templates)
		} else if r.FormValue("source") == "/tasks" {
			task.HandlerCreateTag(w, r, queries, templates)

		}
	})))
	SetupAuthRoutes(mux, queries, templates)
	SetupTaskRoutes(mux, queries, templates)
	SetupTagRoutes(mux, queries, templates)
	SetupVaultRoutes(mux, queries, templates)

	return mux
}
