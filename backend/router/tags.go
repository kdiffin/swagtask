package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"swagtask/database"
	"swagtask/pages"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Tags(e *echo.Echo, dbpool *pgxpool.Pool) {
	e.POST("/tags", func(c echo.Context) error {
		tagValue := c.FormValue("tag")
		sourceValue := c.FormValue("source")

		if sourceValue == "/tasks" {
			return tasksPageTagPostHandler(dbpool, c, tagValue)
		} else if sourceValue == "/tags" {
			return tagsPageTagPostHandler(dbpool, c, tagValue)
		}

		return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))

	})

	e.GET("/tags", func(c echo.Context) error {
		tagsWithTasksOptions, err := database.GetAllTagsWithTasks(dbpool)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		page := pages.NewTagsPage(tagsWithTasksOptions)
		return c.Render(200, "tags-page", page)
	})

	e.DELETE("/tags/:id", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		if errConv != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		_, errRelations := dbpool.Exec(context.Background(), "DELETE FROM tag_task_relations WHERE tag_id = $1", id)
		_, err := dbpool.Exec(context.Background(), "DELETE FROM tags WHERE id = $1", id)

		if err != nil || errRelations != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		return c.NoContent(200)
	})

	e.PUT("/tags/:id", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		if errConv != nil {
			fmt.Println("error here id")
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		_, err := dbpool.Exec(context.Background(), "UPDATE tags SET name=$1 WHERE id = $2", c.FormValue("name"), id)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		tagWithTasks, errTag := database.GetTagWithTasks(dbpool, id)
		if errTag != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return c.Render(200, "tag-card", tagWithTasks)
	})

	e.POST("/tags/:id/tasks", func(c echo.Context) error {
		// TODO: switched from returning on the insert to an exec and then a full fetch
		// bit less performant but more readable, honestly i should switch to sqlc this boilerplate is getting hectic
		id, errConv := getIdAsStr(c)
		taskIdStr := c.FormValue("task_id")
		taskId, errConvTag := strconv.Atoi(taskIdStr)

		if errConvTag != nil || errConv != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		// add task relation
		_, err := dbpool.Exec(context.Background(), "INSERT INTO tag_task_relations (tag_id, task_id) VALUES($1, $2)", id, taskId)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		// get updated task
		tagWithTasks, errTags := database.GetTagWithTasks(dbpool, id)
		if errTags != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		return c.Render(200, "tag-card", tagWithTasks)
	})

	e.DELETE("/tags/:id/tasks", func(c echo.Context) error {
		id, errConv := getIdAsStr(c)
		tagIdStr := c.FormValue("task_id")
		taskId, errConvTag := strconv.Atoi(tagIdStr)

		if errConvTag != nil || errConv != nil {
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		// delete task relation
		_, err := dbpool.Exec(context.Background(), "DELETE FROM tag_task_relations WHERE tag_id = $1 AND task_id = $2", id, taskId)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		// get updated task
		tagWithTasks, errTags := database.GetTagWithTasks(dbpool, id)
		if errTags != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		return c.Render(200, "tag-card", tagWithTasks)
	})

}

func tasksPageTagPostHandler(dbpool *pgxpool.Pool, c echo.Context, tagValue string) error {
	_, err := dbpool.Exec(context.Background(), "INSERT INTO tags (name) VALUES($1)", tagValue)
	if err != nil {
		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	// invalidate cache
	tasksWithTags, errTasks := database.GetAllTasksWithTags(dbpool)
	if errTasks != nil {
		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	page := pages.NewTasksPage(tasksWithTags)

	return c.Render(200, "tasks-container", page)
}

// TODO: refactor this to just send a single tagWithTasks back
func tagsPageTagPostHandler(dbpool *pgxpool.Pool, c echo.Context, tagValue string) error {
	_, err := dbpool.Exec(context.Background(), "INSERT INTO tags (name) VALUES($1)", tagValue)
	if err != nil {
		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	// invalidate cache
	tagsWithTasks, errTags := database.GetAllTagsWithTasks(dbpool)
	if errTags != nil {
		return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.Render(200, "tags-list-container", tagsWithTasks)
}
