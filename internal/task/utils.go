package task

import (
	"net/http"
	"strings"
	db "swagtask/internal/db/generated"
)

func getTaskAvailableTags(allTags []db.Tag, relatedTags []relatedTag) []availableTag {
	// think of this as a set checking if the tag is a tag of the task
	// int is id
	tagExists := make(map[string]bool)
	for _, tag := range relatedTags {
		tagExists[tag.ID] = true
	}

	availableTags := []availableTag{}
	for _, tag := range allTags {
		if !tagExists[tag.ID.String()] {
			availableTags = append(availableTags, availableTag{
				Name: tag.Name,
				ID:   tag.ID.String(),
			})
		}
	}

	return availableTags
}

func filterParams(r *http.Request) TasksPageFilters {
	tag := strings.TrimSuffix(r.URL.Query().Get("tags"), "/")
	task := strings.TrimSuffix(r.URL.Query().Get("taskName"), "/")
	filters := newTasksPageFilters(tag, task)
	return filters
}
