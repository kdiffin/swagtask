package vault

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

func HandlerGetVaults(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vaults, err := getVaultsWithCollaborators(queries, utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	page := newVaultsPage(vaults, true, user.PathToPfp, user.Username)
	templates.Render(w, "vaults-page", page)
}

func HandlerCreateVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vaults, err := createvault(queries, r.FormValue("vault_name"), r.FormValue("vault_description"), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "vault-card", vaults)
}

func HandlerDeleteVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	errDelete := deletevault(queries, utils.PgUUID(r.PathValue("vaultId")), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, errDelete) {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(nil))
}

func HandlerUpdateVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	errDelete := updateVault(queries, utils.PgUUID(r.PathValue("vaultId")), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, errDelete) {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(nil))
}
