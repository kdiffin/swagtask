package vault

import (
	"fmt"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/task"
	"swagtask/internal/template"
	"swagtask/internal/utils"
	common "swagtask/internal/vault/common"
)

func HandlerGetTasks(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	filters := task.FilterParams(r)
	tasks, err := task.GetFilteredTasksWithTags(queries, filters, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
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

	page := common.NewVaultTasksPage(
		tasks,
		filters,
		true,
		user.PathToPfp,
		vaultWithCollaborators.VaultUI,
		common.UserVaultUI{
			PathToPfp: user.PathToPfp,
			Username:  user.Username,
			Role:      role,
		},
		collaborators,
		user.Username)

	templates.Render(w, "collaborative-tasks-page", page)
}
