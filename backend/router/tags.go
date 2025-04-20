package router

import (
	"context"
	"net/http"
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
			return c.String(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
		}

		var tag database.Tag
		err := dbpool.QueryRow(context.Background(), "UPDATE tags SET name=$1 WHERE id = $2 RETURNING id,name", c.FormValue("name"), id).Scan(&tag.Id, &tag.Name)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		rows, errTasks := dbpool.Query(context.Background(), "SELECT name,id FROM tasks")
		if errTasks != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		allTasks := []database.TagRelationOption{}
		for rows.Next() {
			var option database.TagRelationOption
			rows.Scan(&option.Name, &option.Id)

			allTasks = append(allTasks, option)
		}

		tagWithTasks := database.NewTagWithTasks(tag.Id, tag.Name, allTasks)
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
