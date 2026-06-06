package vault

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
	common "swagtask/internal/vault/common"
)

type VaultHandler struct {
	queries   *db.Queries
	templates *template.Template
}

func NewVaultHandler(queries *db.Queries, templates *template.Template) *VaultHandler {
	return &VaultHandler{queries: queries, templates: templates}
}

func (h *VaultHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vaults, err := common.GetVaultsWithCollaborators(h.queries, utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	page := common.NewVaultsPage(vaults, true, user.PathToPfp, user.Username)
	h.templates.Render(w, "vaults-page", page)
}

func (h *VaultHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vaults, err := common.CreateVault(h.queries, r.FormValue("vault_name"), r.FormValue("vault_description"), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	h.templates.Render(w, "vault-card", vaults)
}

func (h *VaultHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	errDelete := common.DeleteVault(h.queries, utils.PgUUID(r.PathValue("vaultId")), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, errDelete) {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(nil))
}

func (h *VaultHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	var locked bool
	if r.FormValue("vault_locked") == "" {
		locked = false
	} else {
		locked = true
	}

	vault, errUpdate := common.UpdateVault(h.queries, utils.PgUUID(r.PathValue("vaultId")), utils.PgUUID(user.ID), r.FormValue("vault_name"), r.FormValue("vault_description"), locked, r.Context())
	if utils.CheckError(w, r, errUpdate) {
		return
	}

	h.templates.Render(w, "vault-card", vault)
}

func (h *VaultHandler) AddCollaborator(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	vault, errUpdate := common.AddCollaboratorToVault(h.queries,
		r.FormValue("collaborator_username"),
		utils.PgUUID(user.ID),
		utils.PgUUID(r.PathValue("vaultId")),
		r.FormValue("collaborator_role"),
		r.Context())
	if utils.CheckError(w, r, errUpdate) {
		return
	}

	h.templates.Render(w, "vault-card", vault)
}

func (h *VaultHandler) RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	collaboratorUsername := r.FormValue("collaborator_username")
	if collaboratorUsername == user.Username {
		http.Error(w, "You cant remove the owner as a collaborator, consider deleting the vault if needed.", http.StatusBadRequest)
		return
	}
	vault, errUpdate := common.RemoveCollaboratorFromVault(h.queries,
		utils.PgUUID(user.ID),
		utils.PgUUID(r.PathValue("vaultId")),
		collaboratorUsername,
		r.Context())
	if utils.CheckError(w, r, errUpdate) {
		return
	}

	h.templates.Render(w, "vault-card", vault)
}
