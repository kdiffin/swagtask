package tag

import (
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

type TagHandler struct {
	queries   *db.Queries
	templates *template.Template
}

func NewTagHandler(queries *db.Queries, templates *template.Template) *TagHandler {
	return &TagHandler{queries: queries, templates: templates}
}

func (h *TagHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}
	tagsWithTasks, errTag := GetTagsWithTasks(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errTag) {
		return
	}

	page := NewTagsPage(tagsWithTasks, true, user.PathToPfp, user.Username)
	h.templates.Render(w, "tags-page", page)
}

func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	tagWithTask, err := UpdateTag(h.queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), tagName, r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	h.templates.Render(w, "tag-card", tagWithTask)
}

func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tagId := r.PathValue("id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	errTag := DeleteTag(h.queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, errTag) {
		return
	}

	// bc we want htmx to rerender it
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nil))
}

func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	tagWithTasks, errTag := CreateTag(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.FormValue("tag_name"), r.Context())
	if utils.CheckError(w, r, errTag) {
		return
	}

	h.templates.Render(w, "tag-card", tagWithTasks)
}

func (h *TagHandler) AddTask(w http.ResponseWriter, r *http.Request) {
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

	tagWithTasks, err := AddTaskToTag(h.queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(taskId), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	h.templates.Render(w, "tag-card", tagWithTasks)
}

func (h *TagHandler) RemoveTask(w http.ResponseWriter, r *http.Request) {
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
	tagWithTasks, err := DeleteTaskRelationFromTag(h.queries, utils.PgUUID(tagId), utils.PgUUID(taskId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	h.templates.Render(w, "tag-card", tagWithTasks)
}
