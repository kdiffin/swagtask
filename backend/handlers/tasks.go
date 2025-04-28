package handlers

import (
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
	"swagtask/utils"
)

// ---- READ ----

func HandlerGetTasks(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template)   {
	tasks, err := service.GetAllTasksWithTags(queries)

	if err != nil {
		utils.LogError("getting all tasks with tags", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	page := models.NewTasksPage(tasks)
	templates.Render(w, "tasks-page", page)
}


// ---- UPDATE ----

func HandlerCreateTask(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template)   {
	name := r.FormValue("name")
	idea := r.FormValue("idea")

	task, err := service.CreateTask(queries, name, idea)

	if err != nil {
		utils.LogError("Creating Task", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		templates.Render(w, "form-error", "DONT ADD A TASK WITH THE SAME IDEA BRO.")
		return
	}

	templates.Render(w, "form-success", nil)
	templates.Render(w, "task", task)
}

func HandlerTaskToggleComplete(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template, taskId int32)   {
	taskWithTags, err := service.UpdateTaskCompletion(queries, taskId)

	if err != nil {
		utils.LogError("Updating Task completion", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerAddTagToTask(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template, taskId int32, tagId int32) {
	taskWithTags, err := service.AddTagToTask(queries, tagId, taskId)
	if err != nil {
		utils.LogError("Updating Task completion", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	templates.Render(w, "task", taskWithTags)
}


//  ---- DELETE ----
func HandlerDeleteTask(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template, taskId int32) {
	err := 	service.DeleteTask(queries, taskId)
	if err != nil {
		utils.LogError("Updating Task completion", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(200)
	w.Write([]byte(nil))
}
	