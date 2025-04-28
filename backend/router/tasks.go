package router

import (
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
)

func Tasks(mux *http.ServeMux, queries *db.Queries, templates *models.Template) {


	
	// e.DELETE("/tasks/:id/tags", func(c echo.Context) error {
	// 	id, errConv := getIdAsStr(c)
	// 	tagIdStr := c.FormValue("tag")
	// 	tagId, errConvTag := strconv.Atoi(tagIdStr)

	// 	if errConvTag != nil {
	// 		return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
	// 	}
	// 	if errConv != nil {
	// 		return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
	// 	}

	// 	// delete task relation
	// 	_, err := dbpool.Exec(context.Background(), "DELETE FROM tag_task_relations WHERE task_id = $1 AND tag_id = $2", id, tagId)
	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	}

	// 	// get updated task
	// 	taskWithTags, errTasks := database.GetTaskWithTagsById(dbpool, id)
	// 	if errTasks != nil {
	// 		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	}

	// 	return c.Render(200, "task", taskWithTags)
	// })

	// // TODO: refactor this function
	// e.PUT("/tasks/:id", func(c echo.Context) error {
	// 	id, errConv := getIdAsStr(c)
	// 	if errConv != nil {
	// 		return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
	// 	}

	// 	// update the task first (no relations)
	// 	name := c.FormValue("name")
	// 	idea := c.FormValue("idea")
	// 	task, errUpdateTask := database.UpdateTask(dbpool, name, idea, id)
	// 	if errUpdateTask != nil {
	// 		// return no contents
	// 		// if theres no update to tasks skip
	// 		if !errors.Is(errUpdateTask, database.ErrNoUpdateFields) {
	// 			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
	// 		}

	// 	}

	// 	tagsOfTask, errTagsOfTask := database.GetTagsOfTask(dbpool, id)
	// 	if errTagsOfTask != nil {
	// 		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	}

	// 	// add tag logic
	// 	allTags, errAllTags := database.GetAllTags(dbpool)
	// 	if errAllTags != nil {
	// 		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	}
	// 	availableTags := database.GetTaskAvailableTags(allTags, tagsOfTask)
	// 	taskWithTags := database.NewTaskWithTags(*task, tagsOfTask, availableTags)
	// 	// dbpool.Query(,str)
	// 	return c.Render(200, "task", taskWithTags)
	// })

	// e.DELETE("/tasks/:id", func(c echo.Context) error {
	// 	id, errConv := getIdAsStr(c)
	// 	if errConv != nil {
	// 		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 	}

	// 	_, errRelations := dbpool.Exec(context.Background(), "DELETE FROM tag_task_relations WHERE task_id = $1", id)
	// 	if errRelations != nil {
	// 		return c.String(http.StatusInternalServerError, "grug got confused on relation deletion")
	// 	}
	// 	_, err := dbpool.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, "grug couldnt delete id")
	// 	}

	// 	return c.NoContent(200)
	// })
	// e.GET("/tasks/:id", func(c echo.Context) error {
	// 	id, errConv := getIdAsStr(c)
	// 	if errConv != nil {
	// 		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	// 	}

	// 	task, err := database.GetTaskWithTagsById(dbpool, id)
	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	}

	// 	var prevButton pages.TaskButton
	// 	var nextButton pages.TaskButton
	// 	prevButton.Exists = true
	// 	nextButton.Exists = true
	// 	errPrevIdExists := dbpool.QueryRow(context.Background(), "SELECT id, name FROM tasks WHERE id < $1 ORDER BY id DESC LIMIT 1", id).Scan(&prevButton.ID, &prevButton.Name)
	// 	errNextIdExists := dbpool.QueryRow(context.Background(), "SELECT id, name FROM tasks WHERE id > $1 ORDER BY id ASC LIMIT 1", id).Scan(&nextButton.ID, &nextButton.Name)
	// 	if errPrevIdExists != nil {
	// 		prevButton.Exists = false
	// 	}
	// 	if errNextIdExists != nil {
	// 		nextButton.Exists = false
	// 	}

	// 	page := pages.NewTaskPage(
	// 		*task,
	// 		prevButton,
	// 		nextButton,
	// 	)
	// 	return c.Render(200, "task-page", page)
	// })

}
