package handlers

import (
	"errors"
	"net/http"
	"strings"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
	"swagtask/utils"
)

// ---- READ ----

func HandlerGetTasks(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template)   {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	tag := strings.TrimSuffix(r.URL.Query().Get("tags") ,"/")
	task := strings.TrimSuffix(r.URL.Query().Get("taskName") ,"/")
	filters := models.NewTasksPageFilters(tag, task)
	tasks, err := service.GetFilteredTasksWithTags(queries, &filters, user.ID, r.Context())
	if checkErrors(w, r, err)  {
		return
	}
	
	page := models.NewTasksPage(tasks, &filters.ActiveTag, &filters.SearchQuery, true, user.PathToPfp, user.Username)
	templates.Render(w, "tasks-page", page)
}

func HandlerGetTask(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template, id int32)   { 
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w, r, errAuth) {
		return
	}
	taskWithTags, err := service.GetTaskWithTagsById(queries,user.ID,id,r.Context())
	if checkErrors(w, r, err) {
		return
	}

	prevButton, nextButton := service.GetTaskNavigationButtons(r.Context(), queries, user.ID, id)
	page := models.NewTaskPage(*taskWithTags, prevButton, nextButton, true, user.PathToPfp, user.Username)
	templates.Render(w, "task-page", page)
}

// ---- CREATE ----

func HandlerCreateTask(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template)   {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	name := r.FormValue("task_name")
	idea := r.FormValue("task_idea")
	
	task, err := service.CreateTask(queries, name,user.ID,idea, r.Context())
	if err != nil {
		if errors.Is(err, service.ErrUnprocessable) {
			utils.LogError("error adding task", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			templates.Render(w, "form-error", "DONT ADD THE TASK WITH THE SAME IDEA VRO")
			return
		} else {
			checkErrors(w,r,  err)
			return
		}
	}
	templates.Render(w, "form-success", nil)
	templates.Render(w, "task", task)
}

// ---- UPDATE ----

func HandlerTaskToggleComplete(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template, taskId int32)   {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	taskWithTags, err := service.UpdateTaskCompletion(queries,user.ID, taskId,r.Context())
	if checkErrors(w,r, err) {
		return
	}

	templates.Render(w, "task", taskWithTags)
}

func HandlerUpdateTask(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template,  taskId int32, idea string, name string) {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	taskWithTags, errUpdate := service.UpdateTask(queries, taskId,user.ID,name,idea,r.Context())
	if errUpdate != nil {
		// return no contents
		// if theres no update to tasks skip
		if errors.Is(errUpdate, service.ErrNoUpdateFields) {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(nil)) 
			return
		}  else if errors.Is(errUpdate, service.ErrUnprocessable) {
			templates.Render(w, "tasks-container-error", "Task has same idea: " + idea)
			taskWithTags,_ := service.GetTaskWithTagsById(queries, user.ID,taskId,r.Context())
			templates.Render(w, "task", taskWithTags)
			return
		} else if checkErrors(w,r, errUpdate) {
			return
		}
	}
	
	templates.Render(w, "tasks-container-success", "Successfully updated task: " + taskWithTags.Name)
	templates.Render(w, "task", taskWithTags)
}

//  ---- DELETE ----
func HandlerDeleteTask(w http.ResponseWriter, r *http.Request ,queries *db.Queries, templates *models.Template, taskId int32) {
	user, errAuth := getUserInfoFromSessionId(queries, r)
	if checkErrors(w,r, errAuth)  {
		return
	}
	err := 	service.DeleteTask(queries, taskId, user.ID,r.Context())
	if checkErrors(w,r, err) {
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(nil))
}
	