package router

import (
	"net/http"
	"strconv"
	"swagtask/database"
	"swagtask/pages"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Tasks(e *echo.Echo, dbpool *pgxpool.Pool) {
	// e.POST("/tasks", func(c echo.Context) error {
	// 	var task database.Task
	// 	err := dbpool.QueryRow(context.Background(), "INSERT INTO tasks (name, idea) VALUES ($1, $2) RETURNING (name, idea, id, tags, completed)",
	// 		c.FormValue("name"),
	// 		c.FormValue("idea")).Scan(&task)

	// 	if err != nil {
	// 		return c.Render(422, "form-error", "U CANT PUT THE SAME IDEA TWICE")

	// 	}

	// 	c.Render(200, "form-success", nil)
	// 	return c.Render(200, "task", task)
	// })

	e.GET("/tasks/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, errConv := strconv.Atoi(idStr)
		if errConv != nil {
			return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}
		task, err := database.GetTaskWithTagsById(dbpool, id)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		page := pages.TaskPage{
			Task:           *task,
			PrevTaskExists: true,
			NextTaskExists: true,
			PrevId:         1,
			NextId:         1,
		}
		return c.Render(200, "task-page", page)
	})

	// e.POST("/tasks/:id/toggle-complete", func(c echo.Context) error {
	// 	idStr := c.Param("id")
	// 	id, errConv := strconv.Atoi(idStr)
	// 	if errConv != nil {
	// 		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadGateway))
	// 	}

	// 	var task database.Task
	// 	err := dbpool.QueryRow(context.Background(), "UPDATE tasks SET completed = NOT completed WHERE id = $1 RETURNING name,idea,id,tags,completed", id).Scan(&task.Name, &task.Idea, &task.Id, &task.Tags, &task.Completed)

	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	}

	// 	return c.Render(200, "task", task)
	// })

	// e.PUT("/tasks/:id", func(c echo.Context) error {
	// 	return editTaskHandler(c, dbpool)
	// })

	// e.DELETE("/tasks/:id", func(c echo.Context) error {
	// 	idStr := c.Param("id")
	// 	id, errConv := strconv.Atoi(idStr)
	// 	if errConv != nil {
	// 		return c.String(http.StatusBadRequest, http.StatusText(http.StatusBadGateway))
	// 	}

	// 	_, err := dbpool.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	// 	}

	// 	return c.NoContent(200)
	// })
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
