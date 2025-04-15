package router

import (
	"context"
	"log"
	"myapp/database"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type TagsPage struct {
	Tasks     database.Tasks
	Tag       string
	OtherTags []string
}

func Tags(e *echo.Echo, dbpool *pgxpool.Pool) {
	e.GET("/tags/:tag", func(c echo.Context) error {
		tag := c.Param("tag")

		rows, err := dbpool.Query(context.Background(), "SELECT id,name,idea,tags,completed FROM tasks WHERE $1 = ANY(tags)", tag)
		rowsOtherTags, errOtherTags := dbpool.Query(context.Background(), "SELECT tags FROM tasks WHERE $1 NOT = ANY(tags)", tag)

		if errOtherTags != nil {
			// If there's an error executing the query, log it and handle gracefully
			log.Printf("Error executing query: %v", errOtherTags)
		}
		if err != nil {
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		var tasks database.Tasks
		for rows.Next() {
			var task database.Task
			err := rows.Scan(&task.Id, &task.Name, &task.Idea, &task.Tags, &task.Completed)
			if err != nil {
				c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

			tasks = append(tasks, task)
		}

		var otherTags []string
		for rowsOtherTags.Next() {
			var tag string
			err := rows.Scan(&tag)

			// gracefully hanndle the error
			if err != nil {
				log.Printf("Error scanning row: %v", err)
			}

			otherTags = append(otherTags, tag)
		}

		page := TagsPage{
			Tasks:     tasks,
			Tag:       tag,
			OtherTags: otherTags,
		}
		return c.Render(200, "tags-page", page)
	})
}

// func Tags(e *echo.Echo, db *pgxpool.Pool) {
// 	e.GET("/tags/:tag", func(c echo.Context) error {
// 		tag := c.Param("tag")

// 		// set which maps the key to the value (both are the same thing cuz this is a set)
// 		tags := make(map[string]string)

// 		var tasksWithTag database.Tasks
// 		for _, task := range db.Tasks {
// 			tagInTask := false

// 			// add tag so it can be passed onto html

// 			for _, tagOfTask := range task.Tags {
// 				tags[tagOfTask] = tagOfTask
// 				if tagOfTask == tag {
// 					tagInTask = true
// 					delete(tags, tagOfTask)
// 				}

// 			}

// 			if tagInTask {
// 				tasksWithTag = append(tasksWithTag, task)
// 			}
// 		}

// 		var otherTags []string
// 		for _, value := range tags {
// 			otherTags = append(otherTags, value)
// 		}

// 		taskPage := TagsPage{
// 			Tasks:     tasksWithTag,
// 			OtherTags: otherTags,
// 			Tag:       tag,
// 		}
// 		return c.Render(200, "tags-page", taskPage)
// 	})
// }
