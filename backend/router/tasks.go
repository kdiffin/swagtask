package router

import (
	"myapp/database"
	"myapp/utils"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type TaskPage struct {
	Task           database.Task
	PrevTaskExists bool
	NextTaskExists bool
	PrevId         int
	NextId         int
}

func Tasks(e *echo.Echo, db *database.Database) {
	e.DELETE("/tasks/:id", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		deleted := false
		for i, task := range db.Tasks {
			if task.Id == id {
				db.Tasks = append(db.Tasks[:i], db.Tasks[i+1:]...)
				deleted = true
				break
			}
		}
		if !deleted {
			return c.String(400, "Task not found")
		}

		return c.NoContent(200)
	})

	e.PUT("/tasks/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		updated := false
		index := -1

		for i, task := range db.Tasks {
			if task.Id == id {
				new_tags := task.Tags
				if c.FormValue("tag") != "" {
					new_tags = append(new_tags, c.FormValue("tag"))
				}

				updatedTask := database.Task{
					Name: utils.StringWithFallback(c.FormValue("name"), task.Name),
					Idea: utils.StringWithFallback(c.FormValue("idea"), task.Idea),
					Id:   id,
					Tags: new_tags,
				}

				db.Tasks[i] = updatedTask
				index = i
				updated = true
				break
			}
		}
		if !updated {
			return c.String(400, "Task not found")
		}

		return c.Render(200, "task", db.Tasks[index])
	})

	e.POST("/tasks/:id/toggle-complete", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		index := -1
		for i, task := range db.Tasks {
			if task.Id == id {
				db.Tasks[i].Completed = !task.Completed
				index = i
			}
		}

		return c.Render(200, "task", db.Tasks[index])
	})

	// default id
	id := 3
	e.POST("/tasks", func(c echo.Context) error {
		name := c.FormValue("name")
		idea := c.FormValue("idea")

		task := database.Task{
			Name: name,
			Idea: idea,
			Id:   id,
		}

		if db.HasIdea(idea) {
			blockData := struct{ ErrorText string }{ErrorText: "STOP ERRORING"}
			return c.Render(422, "form-error", blockData)
		} else {
			c.Render(200, "form-success", "")
		}

		db.Tasks = append([]database.Task{database.NewTask(name, idea, id, []string{})}, db.Tasks...)
		id++

		return c.Render(200, "task", task)
	})

	e.GET("/task/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		var taskOfId TaskPage
		for _, task := range db.Tasks {
			if task.Id == id {
				var pexists, nexists bool

				if task.Id > 1 {
					pexists = true
				}
				if task.Id < len(db.Tasks) {
					nexists = true
				}

				taskOfId = TaskPage{Task: task, PrevTaskExists: pexists, NextTaskExists: nexists, PrevId: task.Id - 1, NextId: task.Id + 1}
			}
		}
		return c.Render(200, "tasks-page", taskOfId)
	})
}
