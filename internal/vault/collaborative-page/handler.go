package vault

import (
	"fmt"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
	common "swagtask/internal/vault/common"
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

	vaultWithCollaborators, errVault := common.GetVaultWithCollaboratorsById(queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errVault) {
		return
	}

	var role string
	collaborators := []common.CollaboratorUI{}
	for _, item := range vaultWithCollaborators.RelatedCollaborators {
		fmt.Println(item.Name)
		if item.Name == user.Username {
			role = item.Role
		}

		collaborators = append(collaborators, common.CollaboratorUI(item))
	}

	page := common.NewVaultPage(
		common.UserVaultUI{
			PathToPfp:  user.PathToPfp,
			Username:   user.Username,
			Authorized: true,
			Role:       role,
		},
		vaultWithCollaborators.VaultUI,
		collaborators)
	templates.Render(w, "vault-page", page)

}
