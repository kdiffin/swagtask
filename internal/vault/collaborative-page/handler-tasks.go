package vault

import (
	"net/http"
	"swagtask/internal/middleware"
	"swagtask/internal/task"
	"swagtask/internal/utils"
	common "swagtask/internal/vault/common"
)

func (h *VaultHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	filters := task.FilterParams(r)
	tasks, err := task.GetFilteredTasksWithTags(h.queries, filters, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	tasksReal := addVaultIDToTasks(vaultId, tasks)

	vaultWithCollaborators, errVault := common.GetVaultWithCollaboratorsById(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errVault) {
		return
	}

	role, collaborators := collaboratorView(vaultWithCollaborators.RelatedCollaborators, user.Username)

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

	h.templates.Render(w, "collaborative-tasks-page", page)
}

// this is for the single page
// i had to replicate this function because of the crap with templates
// not being smart
func (h *VaultHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}
	t, createdAt, err := task.GetTaskPage(h.queries,
		utils.PgUUID(user.ID),
		utils.PgUUID(vaultId),
		utils.PgUUID(id), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	taskWithVault := addVaultIDToTask(vaultId, *t)
	vaultWithCollaborators, errVault := common.GetVaultWithCollaboratorsById(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errVault) {
		return
	}

	role, collaborators := collaboratorView(vaultWithCollaborators.RelatedCollaborators, user.Username)
	prevButton, nextButton := task.GetTaskNavigationButtons(r.Context(), h.queries, createdAt, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(id))
	page := common.NewVaultTaskPage(taskWithVault,
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

	h.templates.Render(w, "collaborative-task-page", page)
}

func addVaultIDToTasks(vaultID string, tasks []task.TaskWithTags) []task.TaskWithTags {
	result := make([]task.TaskWithTags, 0, len(tasks))
	for _, item := range tasks {
		result = append(result, addVaultIDToTask(vaultID, item))
	}
	return result
}

func addVaultIDToTask(vaultId string, t task.TaskWithTags) task.TaskWithTags {
	tasksReal := task.TaskWithTags{
		TaskUI: task.TaskUI{
			CreatedAt: t.CreatedAt,
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
	return tasksReal
}
