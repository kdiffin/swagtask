package vault

import (
	"fmt"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	"swagtask/internal/template"
	common "swagtask/internal/vault/common"

	"swagtask/internal/utils"
)

func HandlerGetTags(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	tagsWithTasks, errTags := tag.GetTagsWithTasks(queries, utils.PgUUID(user.ID), utils.PgUUID(user.DefaultVaultID), r.Context())
	if utils.CheckError(w, r, errTags) {
		fmt.Println("error was here")
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

	page := common.NewVaultTagsPage(tagsWithTasks,
		vaultWithCollaborators.VaultUI,
		common.UserVaultUI{
			PathToPfp:  user.PathToPfp,
			Authorized: true,
			Username:   user.Username,
			Role:       role,
		},
		collaborators)
	templates.Render(w, "collaborative-tags-page", page)

}
