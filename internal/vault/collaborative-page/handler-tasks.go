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

	// AGAIN THIS HACKY BULLSHIT BC OF NOT USING TEMPL :SOB
	// FEELING THE TECH DEBT
	var tasksReal []task.TaskWithTags
	for _, t := range tasks {
		tasksReal = append(tasksReal, task.TaskWithTags{
			TaskUI: task.TaskUI{
				ID:        t.ID,
				Name:      t.Name,
				Author:    t.Author,
				Idea:      t.Idea,
				Completed: t.Completed,
			},
			VaultID:       vaultId,
			RelatedTags:   t.RelatedTags,
			AvailableTags: t.AvailableTags,
		})
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
		tasksReal,
		filters,
		vaultWithCollaborators.VaultUI,
		common.UserVaultUI{
			PathToPfp:  user.PathToPfp,
			Authorized: true,
			Username:   user.Username,
			Role:       role,
		},
		collaborators,
	)

	templates.Render(w, "collaborative-tasks-page", page)
}

// this is for the single page
// i had to replicate this function because of the crap with templates
// not being smart
func HandlerGetTask(w http.ResponseWriter,
	r *http.Request,
	queries *db.Queries,

	templates *template.Template) {
	id := r.PathValue("id")

	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}
	t, createdAt, err := task.GetTaskPage(queries,
		utils.PgUUID(user.ID),
		utils.PgUUID(vaultId),
		utils.PgUUID(id), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	tasksReal := task.TaskWithTags{
		TaskUI: task.TaskUI{
			ID:        t.ID,
			Name:      t.Name,
			Author:    t.Author,
			Idea:      t.Idea,
			Completed: t.Completed,
		},
		VaultID:       vaultId,
		RelatedTags:   t.RelatedTags,
		AvailableTags: t.AvailableTags,
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
	prevButton, nextButton := task.GetTaskNavigationButtons(r.Context(), queries, createdAt, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(id))
	page := common.NewVaultTaskPage(tasksReal,
		task.TaskPageButtons{
			PrevButton: prevButton,
			NextButton: nextButton,
		},
		vaultWithCollaborators.VaultUI,
		common.UserVaultUI{
			PathToPfp:  user.PathToPfp,
			Username:   user.Username,
			Authorized: true,
			Role:       role,
		},
		collaborators,
	)

	templates.Render(w, "collaborative-task-page", page)
}
