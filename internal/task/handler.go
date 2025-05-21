package task

import (
	"errors"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

// ---- READ ----

func HandlerAddTagToTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
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

	taskWithTags, err := addTagToTask(queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(taskId), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerRemoveTagFromTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
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
	taskWithTags, err := deleteTagRelationFromTask(queries, utils.PgUUID(tagId), utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(taskId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerGetTasks(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	filters := filterParams(r)
	tasks, err := getFilteredTasksWithTags(queries, filters, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	page := newTasksPage(tasks, filters, true, user.PathToPfp, user.Username)
	templates.Render(w, "tasks-page", page)
}

func HandlerGetTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	id := r.PathValue("id")

	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}
	taskWithTags, createdAt, err := getTaskPage(queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(id), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	prevButton, nextButton := getTaskNavigationButtons(r.Context(), queries, createdAt, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(id))
	page := newTaskPage(*taskWithTags, prevButton, nextButton, true, user.PathToPfp, user.Username)
	templates.Render(w, "task-page", page)
}

// ---- CREATE ----

func HandlerCreateTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
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
	filters := filterParams(r)

	task, err := createTask(queries, name, idea, filters, utils.PgUUID(user.ID), utils.PgUUID(vaultId), r.Context())
	if err != nil {
		if errors.Is(err, utils.ErrUnprocessable) {
			utils.LogError("error adding task", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			templates.Render(w, "form-error", "DONT ADD THE TASK WITH THE SAME IDEA VRO")
			return
		} else {
			utils.CheckError(w, r, err)
			return
		}
	}

	templates.Render(w, "form-success", nil)
	templates.Render(w, "task", task)
}

// ---- UPDATE ----

func HandlerTaskToggleComplete(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	taskId := r.PathValue("id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	taskWithTags, err := updateTaskCompletion(queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(taskId), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerUpdateTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
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

	taskWithTags, errUpdate := updateTask(queries, utils.PgUUID(vaultId), utils.PgUUID(taskId), utils.PgUUID(user.ID), name, idea, r.Context())
	if errUpdate != nil {
		// return no contents
		// if theres no update to tasks skip
		if errors.Is(errUpdate, utils.ErrNoUpdateFields) {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(nil))
			return
		} else if errors.Is(errUpdate, utils.ErrUnprocessable) {
			templates.Render(w, "tasks-container-error", "Task has same idea: "+idea)
			taskWithTags, _ := getTaskWithTagsById(queries, utils.PgUUID(user.ID), utils.PgUUID(vaultId), utils.PgUUID(taskId), r.Context())
			templates.Render(w, "task", taskWithTags)
			return
		} else if utils.CheckError(w, r, errUpdate) {
			return
		}
	}

	templates.Render(w, "tasks-container-success", "Successfully updated task: "+taskWithTags.Name)
	templates.Render(w, "task", taskWithTags)
}

// ---- DELETE ----
func HandlerDeleteTask(w http.ResponseWriter, r *http.Request, queries *db.Queries, templates *template.Template) {
	taskId := r.PathValue("id")
	user, err := middleware.UserFromContext(r.Context())
	vaultId, errVault := middleware.VaultIDFromContext(r.Context())
	if utils.CheckError(w, r, err) {
		return
	}
	if utils.CheckError(w, r, errVault) {
		return
	}

	err = deleteTask(queries, utils.PgUUID(taskId), utils.PgUUID(vaultId), utils.PgUUID(user.ID), r.Context())
	if utils.CheckError(w, r, err) {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(nil))
}
