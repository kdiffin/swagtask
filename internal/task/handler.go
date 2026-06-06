package task

import (
	"errors"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

type TaskHandler struct {
	queries   *db.Queries
	templates *template.Template
}

func NewTaskHandler(queries *db.Queries, templates *template.Template) *TaskHandler {
	return &TaskHandler{queries: queries, templates: templates}
}

// ---- READ ----

func (h *TaskHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	taskId := r.PathValue("id")
	tagId := r.FormValue("tag_id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	taskWithTags, err := AddTagToTask(h.queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(taskId), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	h.templates.Render(w, "task", taskWithTags)
}

func (h *TaskHandler) CreateTagOption(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	_, errTag := h.queries.CreateTag(r.Context(), db.CreateTagParams{
		Name:    r.FormValue("tag_name"),
		UserID:  utils.PgUUID(user.ID),
		VaultID: utils.PgUUID(vaultId),
	})
	if utils.CheckError(w, r, errTag) {
		return
	}

	filters := FilterParams(r)
	tasks, err := GetFilteredTasksWithTags(h.queries, filters, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	h.templates.Render(w, "tasks-container", tasks)
}

func (h *TaskHandler) RemoveTag(w http.ResponseWriter, r *http.Request) {
	taskId := r.PathValue("id")
	tagId := r.FormValue("tag_id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	taskWithTags, err := DeleteTagRelationFromTask(h.queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(taskId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	h.templates.Render(w, "task", taskWithTags)
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	filters := FilterParams(r)
	tasks, err := GetFilteredTasksWithTags(h.queries, filters, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	page := newTasksPage(tasks, filters, true, user.PathToPfp, user.Username)
	h.templates.Render(w, "tasks-page", page)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}
	taskWithTags, createdAt, err := GetTaskPage(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(id), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	prevButton, nextButton := GetTaskNavigationButtons(r.Context(), h.queries, createdAt, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(id))
	page := NewTaskPage(*taskWithTags, prevButton, nextButton, true, user.PathToPfp, user.Username)
	h.templates.Render(w, "task-page", page)
}

// ---- CREATE ----

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}
	name := r.FormValue("task_name")
	idea := r.FormValue("task_idea")

	task, err := CreateTask(h.queries, name, idea, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if err != nil {
		if errors.Is(err, utils.ErrUnprocessable) {
			utils.LogError("error adding task", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			h.templates.Render(w, "form-error", "DONT ADD THE TASK WITH THE SAME IDEA VRO")
			return
		} else {
			utils.CheckError(w, r, err)
			return
		}
	}

	h.templates.Render(w, "form-success", nil)
	h.templates.Render(w, "task", task)
}

// ---- UPDATE ----

func (h *TaskHandler) ToggleComplete(w http.ResponseWriter, r *http.Request) {
	taskId := r.PathValue("id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	taskWithTags, err := UpdateTaskCompletion(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(taskId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	h.templates.Render(w, "task", taskWithTags)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskId := r.PathValue("id")
	name := r.FormValue("task_name")
	idea := r.FormValue("task_idea")

	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	taskWithTags, errUpdate := UpdateTask(h.queries, utils.PgUUID(vaultId), utils.PgUUID(taskId), utils.PgUUID(user.ID), name, idea, r.Context())
	if errUpdate != nil {
		// return no contents
		// if theres no update to tasks skip
		if errors.Is(errUpdate, utils.ErrNoUpdateFields) {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(nil))
			return
		} else if errors.Is(errUpdate, utils.ErrUnprocessable) {
			h.templates.Render(w, "tasks-container-error", "Task has same idea: "+idea)
			taskWithTags, _ := GetTaskWithTagsById(h.queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(taskId), r.Context())
			h.templates.Render(w, "task", taskWithTags)
			return
		} else if utils.CheckError(w, r, errUpdate) {
			return
		}
	}

	h.templates.Render(w, "tasks-container-success", "Successfully updated task: "+taskWithTags.Name)
	h.templates.Render(w, "task", taskWithTags)
}

// ---- DELETE ----
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskId := r.PathValue("id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	err = DeleteTask(h.queries, utils.PgUUID(taskId), utils.PgUUID(vaultId), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(nil))
}
