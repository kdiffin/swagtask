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

func (h *VaultHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, errUser := middleware.UserFromContext(r.Context())
	if utils.CheckError(w, r, errUser) {
		return
	}
	vaultId, errVaultId := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, errVaultId) {
		return
	}

	vaultWithCollaborators, errVault := common.GetVaultWithCollaboratorsById(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errVault) {
		return
	}

	role, collaborators := collaboratorView(vaultWithCollaborators.RelatedCollaborators, user.Username)

	page := common.NewVaultPage(
		common.UserVaultUI{
			PathToPfp:  user.PathToPfp,
			Username:   user.Username,
			Authorized: true,
			Role:       role,
		},
		vaultWithCollaborators.VaultUI,
		collaborators)
	h.templates.Render(w, "vault-page", page)

}

func collaboratorView(items []common.RelatedCollaborator, username string) (string, []common.CollaboratorUI) {
	var role string
	collaborators := make([]common.CollaboratorUI, 0, len(items))
	for _, item := range items {
		if item.Name == username {
			role = item.Role
		}
		collaborators = append(collaborators, common.CollaboratorUI(item))
	}
	return role, collaborators
}
