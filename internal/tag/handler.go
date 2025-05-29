package tag

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
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
	tagsWithTasks, errTag := GetTagsWithTasks(queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errTag) {
		return
	}

	page := NewTagsPage(tagsWithTasks, true, user.PathToPfp, user.Username)
	templates.Render(w, "tags-page", page)
}

func HandlerUpdateTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	tagId := r.PathValue("id")
	tagName := r.FormValue("tag_name")

	if tagName == "" {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(nil))
		return
	}
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	tagWithTask, err := updateTag(queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), tagName, r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTask)
}

func HandlerDeleteTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	tagId := r.PathValue("id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	errTag := deleteTag(queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errTag) {
		return
	}

	// bc we want htmx to rerender it
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nil))
}

func HandlerAddTaskToTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	tagId := r.PathValue("id")
	taskId := r.FormValue("task_id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	tagWithTasks, err := addTaskToTag(queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(taskId), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTasks)
}

func HandlerRemoveTaskFromTag(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	tagId := r.PathValue("id")
	taskId := r.FormValue("task_id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}
	tagWithTasks, err := deleteTaskRelationFromTag(queries, utils.PgUUID(tagId), utils.PgUUID(taskId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "tag-card", tagWithTasks)
}
