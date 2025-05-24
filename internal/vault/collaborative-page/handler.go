package vault

import (
	"net/http"
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
	vault "swagtask/internal/vault/owner-dashboard"
)

func HandlerGetVault(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, errUser := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, errUser) {
		return
	}
	vaultId, errVaultId := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, errVaultId) {
		return
	}

	page := newVaultPage(true, userVaultUI{
		PathToPfp: user.PathToPfp,
		Username:  user.Username,
		Role:      "deez",
	}, user.PathToPfp, user.Username, vault.VaultUI{
		ID:   vaultId,
		Name: "best vault",
		Author: auth.Author{
			PathToPfp: "",
			Username:  "grug brained dev",
		},
		Description: "bestest vault ever",
		Locked:      false,
		CreatedAt:   "",
		Kind:        "",
	}, []collaboratorUI{})
	templates.Render(w, "vault-page", page)

}
