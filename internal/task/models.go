package task

import (
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
)

type taskUI struct {
	ID        string
	Name      string
	Author    auth.Author
	Idea      string
	Completed bool
}

func newUITask(task db.Task, author auth.Author) taskUI {

	return taskUI{
		ID:        task.ID.String(),
		Author:    author,
		Name:      task.Name,
		Idea:      task.Idea,
		Completed: task.Completed,
	}
}

// ---- FOR UI ----
// tasks
type tagOption struct {
	Name string
	ID   string
}
type availableTag = tagOption
type relatedTag tagOption
type taskWithTags struct {
	taskUI
	Author        auth.Author
	RelatedTags   []relatedTag
	AvailableTags []availableTag
}

func newTaskWithTags(task taskUI, relatedTags []relatedTag, availableTags []availableTag) taskWithTags {
	return taskWithTags{
		taskUI:        task,
		RelatedTags:   relatedTags,
		AvailableTags: availableTags,
		Author:        task.Author,
	}
}
