package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"swagtask/database"
	"swagtask/pages"
	"swagtask/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// funny right
func getIdAsStr(c echo.Context) (int, error) {
	idStr := c.Param("id")
	id, errConv := strconv.Atoi(idStr)

	return id, errConv
}

func Tasks(e *echo.Echo, dbpool *pgxpool.Pool) {
	e.POST("/tasks", func(c echo.Context) error {
		task, err := database.CreateTask(dbpool, c.FormValue("name"), c.FormValue("idea"))
		if err != nil {
			return c.Render(422, "form-error", "U CANT PUT THE SAME IDEA TWICE")
		}

		taskWithTag := database.TaskWithTags{
			Task: *task,
			Tags: []database.Tag{},
		}

		c.Render(200, "form-success", nil)
		return c.Render(200, "task", taskWithTag)
	})

	e.POST("/tasks/:id/toggle-complete", func(c echo.Context) error {
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

		// Task with or without tag still render
		taskWithTags := database.TaskWithTags{
			Task: task,
			Tags: tags,
		}
		return c.Render(200, "task", taskWithTags)
	})

	// TODO: refactor this function
	e.PUT("/tasks/:id", func(c echo.Context) error {
		updateTaskString := "UPDATE tasks SET"
		var task database.Task
		tagsOfTask := []database.Tag{}
		args := []interface{}{}
		n := 1

		id, errConv := getIdAsStr(c)
		if errConv != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}
		name := c.FormValue("name")
		idea := c.FormValue("idea")
		tag := c.FormValue("tag")
		updateTags := false

		if name != "" {
			updateTaskString += fmt.Sprintf(" name = $%v,", n)
			args = append(args, name)
			n++
		}
		if idea != "" {
			updateTaskString += fmt.Sprintf(" idea = $%v,", n)
			args = append(args, idea)
			n++
		}
		if tag != "" {
			updateTags = true
		}
		args = append(args, id)

		// remove trailing comma
		updateTaskString = strings.TrimSuffix(updateTaskString, ",")
		str := updateTaskString + fmt.Sprintf(" WHERE id = $%v RETURNING name,idea,id,completed", n)
		errTask := dbpool.QueryRow(context.Background(), str, args...).Scan(&task.Name, &task.Idea, &task.Id, &task.Completed)
		utils.PrintList(args)
		fmt.Println(str)
		if errTask != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		var taskWithTags database.TaskWithTags

		rows, errTags := dbpool.Query(context.Background(), `SELECT tg.name, tg.id FROM tags tg JOIN tag_task_relations rel ON rel.tag_id = tg.id WHERE rel.task_id = $1`, id)
		if errTags != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		for rows.Next() {
			var tag database.Tag
			rows.Scan(&tag.Name, &tag.Id)
			tagsOfTask = append(tagsOfTask, tag)
		}

		// add tag logic
		if updateTags {
			var tagObj database.Tag

			errCreateTag := dbpool.QueryRow(context.Background(), "INSERT INTO tags (name) VALUES($1) RETURNING name,id", tag).Scan(&tagObj.Name, &tagObj.Id)
			if errCreateTag != nil {
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			_, errCreateTagToTaskRelation := dbpool.Exec(context.Background(), "INSERT INTO tag_task_relations (tag_id, task_id) VALUES($1, $2)", tagObj.Id, task.Id)
			if errCreateTagToTaskRelation != nil {
				return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			tagsOfTask = append(tagsOfTask, tagObj)
		}

		taskWithTags = database.TaskWithTags{Task: task, Tags: tagsOfTask}

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

		page := pages.TaskPage{
			Task: *task,
			Buttons: struct {
				PrevButton pages.TaskButton
				NextButton pages.TaskButton
			}{
				PrevButton: prevButton,
				NextButton: nextButton,
			},
		}
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
