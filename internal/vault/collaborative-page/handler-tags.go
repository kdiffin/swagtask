package vault

import (
	"net/http"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	common "swagtask/internal/vault/common"

	"swagtask/internal/utils"
)

func (h *VaultHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId :=  r.PathValue("vaultId")
	if utils.CheckError(w, r, err) {
		return
	}


	tagsWithTasks, errTags := tag.GetTagsWithTasks(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errTags) {
		return
	}

	vaultWithCollaborators, errVault := common.GetVaultWithCollaboratorsById(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errVault) {
		return
	}

	role, collaborators := collaboratorView(vaultWithCollaborators.RelatedCollaborators, user.Username)

	page := common.NewVaultTagsPage(tagsWithTasks,
		vaultWithCollaborators.VaultUI,
		common.UserVaultUI{
			PathToPfp:  user.PathToPfp,
			Authorized: true,
			Username:   user.Username,
			Role:       role,
		},
		collaborators)
	h.templates.Render(w, "collaborative-tags-page", page)

}
