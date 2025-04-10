package router

import (
	task_package "myapp/task"
	"myapp/utils"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)




func Tasks(e *echo.Echo, page *task_package.TasksContainer) {
	e.DELETE("/tasks/:id", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		deleted := false
		for i, task := range page.Tasks {
			if task.Id == id {
				page.Tasks = append(page.Tasks[:i], page.Tasks[i+1:]...)
				deleted = true
				break
			}
		}
		if !deleted {
			return c.String(400, "Task not found")
		}

		return c.NoContent(200)
	})

    e.PATCH("/tasks/:id", func(c echo.Context) error {
		idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)		
        if err != nil {
			return c.String(400, "Id must be an integer")
		}
        
       

		updated := false
        index := 0

		for i, task := range page.Tasks {
			if task.Id == id {
                new_tags := append(task.Tags, c.FormValue("tag"))
                updatedTask := task_package.Task{
                    Name: utils.StringWithFallback(c.FormValue("name"), task.Name),
                    Idea: utils.StringWithFallback(c.FormValue("idea"), task.Idea),
                    Id: id,
                    Tags: new_tags,
                }

				page.Tasks[i] = updatedTask
                index = i
				updated = true
				break
			}
		}
		if !updated {
			return c.String(400, "Task not found")
		}

		return c.Render(200, "task", page.Tasks[index])
	})

	// default id
	id := 3
	e.POST("/tasks", func(c echo.Context) error {
		name := c.FormValue("name")
		idea := c.FormValue("idea")
        
		task := task_package.Task{
			Name: name,
			Idea: idea,
			Id:   id,
		}

		if page.HasIdea(idea) {
			blockData := struct{ ErrorText string }{ErrorText: "STOP ERRORING"}
			return c.Render(422, "form-error", blockData)
		} else {
			c.Render(200, "form-success", "")
		}

        
		page.Tasks = append([]task_package.Task{task_package.NewTask(name, idea, id, []string{})}, page.Tasks...)
		id++

		return c.Render(200, "task", task)
	})
}