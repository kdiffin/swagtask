package task

import db "swagtask/internal/db/generated"

type taskUI struct {
	ID        string
	Name      string
	Idea      string
	Completed bool
}

func newUITask(task db.Task) taskUI {

	return taskUI{
		ID:        task.ID.String(),
		Name:      task.Name,
		Idea:      task.Idea,
		Completed: task.Completed,
	}
}

// ---- FOR UI ----
// tasks
type availableTag struct {
	Name string
	ID   string
}
type relatedTag struct {
	Name string
	ID   string
}
type taskWithTags struct {
	Task          taskUI
	RelatedTags   []relatedTag
	AvailableTags []availableTag
}

func NewTaskWithTags(task taskUI, relatedTags []relatedTag, availableTags []availableTag) taskWithTags {
	return taskWithTags{
		Task:          task,
		RelatedTags:   relatedTags,
		AvailableTags: availableTags,
	}
}
