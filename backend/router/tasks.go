package router

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"swagtask/database"
	"swagtask/pages"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Tasks(e *echo.Echo, dbpool *pgxpool.Pool) {
	e.GET("/tasks", func(c echo.Context) error {
		tagNameParam := c.QueryParam("tags")
		taskNameParam := c.QueryParam("taskName")

		if tagNameParam == "" && taskNameParam == "" {
			taskWithTags, err := database.GetAllTasksWithTags(dbpool)
			if err != nil {
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			page := pages.NewTasksPage(
				taskWithTags,
			)
			return c.Render(200, "tasks-page", page)
		}

		filters := database.NewTasksPageFilters(tagNameParam, taskNameParam)
		taskWithTags, errFilteredTasks := database.GetAllFilteredTasksWithTags(dbpool, filters)
		if errFilteredTasks != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		page := pages.NewTasksPage(taskWithTags)
		return c.Render(200, "tasks-page", page)
	})

	e.POST("/tasks", func(c echo.Context) error {
		task, err := database.CreateTask(dbpool, c.FormValue("name"), c.FormValue("idea"))
		if err != nil {
			return c.Render(http.StatusUnprocessableEntity, "form-error", "U CANT PUT THE SAME IDEA TWICE")
		}

		allTags, errAllTags := database.GetAllTags(dbpool)
		if errAllTags != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		taskWithTag := database.NewTaskWithTags(
			*task,
			[]database.Tag{},
			allTags,
		)

		c.Render(200, "form-success", nil)
		return c.Render(200, "task", taskWithTag)
	})

	e.POST("/tasks/:id/toggle-complete", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		id, errConv := getIdAsStr(c)
		if errConv != nil {
			return c.String(http.StatusBadRequest, "bad id")
		}

		row := dbpool.QueryRow(context.Background(), `
			UPDATE tasks SET completed = NOT completed WHERE id = $1 RETURNING id, name, idea, completed
		`, id)

		var task database.Task
		err := row.Scan(&task.Id, &task.Name, &task.Idea, &task.Completed)
		if err != nil {
			return c.String(http.StatusInternalServerError, "grug no find task to update")
		}

		rows, err := dbpool.Query(context.Background(), `
			SELECT tg.id, tg.name FROM tags tg
			INNER JOIN tag_task_relations rel ON tg.id = rel.tag_id
			WHERE rel.task_id = $1
		`, id)
		if err != nil {
			return c.String(http.StatusInternalServerError, "grug no get tags")
		}

		tags := []database.Tag{}
		for rows.Next() {
			var tag database.Tag
			err := rows.Scan(&tag.Id, &tag.Name)
			if err != nil {
				return c.String(http.StatusInternalServerError, "grug confused by tag row")
			}
			tags = append(tags, tag)
		}

		allTags, errAllTags := database.GetAllTags(dbpool)
		if errAllTags != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		availableTags := database.GetTaskAvailableTags(allTags, tags)
		// Task with or without tag still render
		taskWithTags := database.NewTaskWithTags(
			task,
			tags,
			availableTags,
		)
		return c.Render(200, "task", taskWithTags)
	})

	e.POST("/tasks/:id/tags", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		tagIdStr := c.FormValue("tag")
		tagId, errConvTag := strconv.Atoi(tagIdStr)

		if errConvTag != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}
		if errConv != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		// add task relation
		_, err := dbpool.Exec(context.Background(), "INSERT INTO tag_task_relations (task_id, tag_id) VALUES($1, $2)", id, tagId)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		// get updated task
		taskWithTags, errTasks := database.GetTaskWithTagsById(dbpool, id)
		if errTasks != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		return c.Render(200, "task", taskWithTags)
	})

	e.DELETE("/tasks/:id/tags", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		tagIdStr := c.FormValue("tag")
		tagId, errConvTag := strconv.Atoi(tagIdStr)

		if errConvTag != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}
		if errConv != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		// delete task relation
		_, err := dbpool.Exec(context.Background(), "DELETE FROM tag_task_relations WHERE task_id = $1 AND tag_id = $2", id, tagId)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		// get updated task
		taskWithTags, errTasks := database.GetTaskWithTagsById(dbpool, id)
		if errTasks != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		return c.Render(200, "task", taskWithTags)
	})

	// TODO: refactor this function
	e.PUT("/tasks/:id", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		if errConv != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		// update the task first (no relations)
		name := c.FormValue("name")
		idea := c.FormValue("idea")
		task, errUpdateTask := database.UpdateTask(dbpool, name, idea, id)
		if errUpdateTask != nil {
			// return no contents
			// if theres no update to tasks skip
			if !errors.Is(errUpdateTask, database.ErrNoUpdateFields) {
				return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
			}

		}

		tagsOfTask, errTagsOfTask := database.GetTagsOfTask(dbpool, id)
		if errTagsOfTask != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		// add tag logic
		allTags, errAllTags := database.GetAllTags(dbpool)
		if errAllTags != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		availableTags := database.GetTaskAvailableTags(allTags, tagsOfTask)
		taskWithTags := database.NewTaskWithTags(*task, tagsOfTask, availableTags)
		// dbpool.Query(,str)
		return c.Render(200, "task", taskWithTags)
	})

	e.DELETE("/tasks/:id", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		if errConv != nil {
			return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		_, errRelations := dbpool.Exec(context.Background(), "DELETE FROM tag_task_relations WHERE task_id = $1", id)
		if errRelations != nil {
			return c.String(http.StatusInternalServerError, "grug got confused on relation deletion")
		}
		_, err := dbpool.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
		if err != nil {
			return c.String(http.StatusInternalServerError, "grug couldnt delete id")
		}

		return c.NoContent(200)
	})
	e.GET("/tasks/:id", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		if errConv != nil {
			return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		task, err := database.GetTaskWithTagsById(dbpool, id)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		var prevButton pages.TaskButton
		var nextButton pages.TaskButton
		prevButton.Exists = true
		nextButton.Exists = true
		errPrevIdExists := dbpool.QueryRow(context.Background(), "SELECT id, name FROM tasks WHERE id < $1 ORDER BY id DESC LIMIT 1", id).Scan(&prevButton.Id, &prevButton.Name)
		errNextIdExists := dbpool.QueryRow(context.Background(), "SELECT id, name FROM tasks WHERE id > $1 ORDER BY id ASC LIMIT 1", id).Scan(&nextButton.Id, &nextButton.Name)
		if errPrevIdExists != nil {
			prevButton.Exists = false
		}
		if errNextIdExists != nil {
			nextButton.Exists = false
		}

		page := pages.NewTaskPage(
			*task,
			prevButton,
			nextButton,
		)
		return c.Render(200, "task-page", page)
	})

}

// func getTaskHandler(c echo.Context, dbpool *pgxpool.Pool) error {
// 	idStr := c.Param("id")
// 	id, errConv := strconv.Atoi(idStr)
// 	if errConv != nil {
// 		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
// 	}

// 	var task database.Task
// 	err := dbpool.QueryRow(context.Background(),
// 		`SELECT name FROM tasks t
// 		JOIN tag_task_relations rel ON tasks.id = task_id
// 		JOIN ON tags  = WHERE id = $1`, id).Scan(&task.Name, &task.Idea, &task.Completed, &task.Id)

// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
// 	}

// 	var prevId int
// 	var nextId int
// 	prevIdExists := true
// 	nextIdExists := true

// 	// Previous ID
// 	errPrevId := dbpool.QueryRow(context.Background(),
// 		"SELECT id FROM tasks WHERE id < $1 ORDER BY id DESC LIMIT 1", id,
// 	).Scan(&prevId)
// 	if errPrevId != nil {
// 		prevIdExists = false
// 	}

// 	// Next ID
// 	errNextId := dbpool.QueryRow(context.Background(),
// 		"SELECT id FROM tasks WHERE id > $1 ORDER BY id ASC LIMIT 1", id,
// 	).Scan(&nextId)
// 	if errNextId != nil {
// 		nextIdExists = false
// 	}

// 	page := pages.TaskPage{
// 		Task:           task,
// 		PrevTaskExists: prevIdExists,
// 		NextTaskExists: nextIdExists,
// 		PrevId:         prevId,
// 		NextId:         nextId,
// 	}
// 	return c.Render(200, "task-page", page)
// }

// func editTaskHandler(c echo.Context, dbpool *pgxpool.Pool) error {
// 	// TODO: refactor
// 	idStr := c.Param("id")
// 	id, errConv := strconv.Atoi(idStr)
// 	if errConv != nil {
// 		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
// 	}

// 	// dbpool requires any
// 	var args []any
// 	name := c.FormValue("name")
// 	idea := c.FormValue("idea")
// 	tag := c.FormValue("tag")
// 	queryNovalues := "UPDATE tasks SET"
// 	sum := 1
// 	if name != "" {
// 		queryNovalues += fmt.Sprintf(" name = $%d,", sum)
// 		args = append(args, name)
// 		sum++
// 	}
// 	if idea != "" {
// 		queryNovalues += fmt.Sprintf(" idea = $%d,", sum)
// 		args = append(args, idea)
// 		sum++
// 	}
// 	if tag != "" {
// 		queryNovalues += fmt.Sprintf(" tags = array_append(tags, $%d)", sum)
// 		args = append(args, tag)
// 		sum++
// 	}

// 	if len(args) == 0 {
// 		// no updates
// 		return c.NoContent(http.StatusNoContent)
// 	}
// 	queryNovalues = strings.TrimSuffix(queryNovalues, ",")
// 	queryString := queryNovalues + fmt.Sprintf(" WHERE id = $%d RETURNING (name,idea,id,tags,completed)", sum)
// 	args = append(args, id)
// 	println(queryString)
// 	utils.PrintList(args)
// 	// have updates which dynamically generates the querynovalues,

// 	var task database.Task
// 	err := dbpool.QueryRow(context.Background(), queryString, args...).Scan(&task)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
// 	}

// 	return c.Render(200, "task", task)
// }
