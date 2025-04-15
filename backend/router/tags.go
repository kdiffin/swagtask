package router

import (
	"myapp/database"

	"github.com/labstack/echo/v4"
)

type TagsPage struct {
	Tasks     database.Tasks
	Tag       string
	OtherTags []string
}

func Tags(e *echo.Echo, db *database.Database) {
	e.GET("/tags/:tag", func(c echo.Context) error {
		tag := c.Param("tag")

		// set which maps the key to the value (both are the same thing cuz this is a set)
		tags := make(map[string]string)

		var tasksWithTag database.Tasks
		for _, task := range db.Tasks {
			tagInTask := false

			// add tag so it can be passed onto html

			for _, tagOfTask := range task.Tags {
				tags[tagOfTask] = tagOfTask
				if tagOfTask == tag {
					tagInTask = true
					delete(tags, tagOfTask)
				}

			}

			if tagInTask {
				tasksWithTag = append(tasksWithTag, task)
			}
		}

		var otherTags []string
		for _, value := range tags {
			otherTags = append(otherTags, value)
		}

		taskPage := TagsPage{
			Tasks:     tasksWithTag,
			OtherTags: otherTags,
			Tag:       tag,
		}
		return c.Render(200, "tags-page", taskPage)
	})
}
