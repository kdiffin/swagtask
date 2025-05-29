package task

import (
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
)

type TaskUI struct {
	ID        string
	Name      string
	Author    auth.Author
	Idea      string
	Completed bool
}

func newUITask(task db.Task, author auth.Author) TaskUI {

	return TaskUI{
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
type TaskWithTags struct {
	// THIS IS A STUPID HACK FOR THE VAULTS PAGE
	// I CBA TO COPY THE CODE FOR THE VAULTS PART, SO I JUST ADDED THIS
	// AFTER THIS I REALIZED I SHOULD PROB USE TEMPL OR SOMETHING.
	TaskUI
	VaultID       string
	RelatedTags   []relatedTag
	AvailableTags []availableTag
}

func newTaskWithTags(task TaskUI, relatedTags []relatedTag, availableTags []availableTag) TaskWithTags {
	return TaskWithTags{
		TaskUI:        task,
		RelatedTags:   relatedTags,
		AvailableTags: availableTags,
	}
}
