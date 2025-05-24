package tag

import "swagtask/internal/auth"

type tagAuthor = auth.Author

type tagUI struct {
	VaultID string
	Author  tagAuthor
	Name    string
	ID      string
}
type taskOption struct {
	Name string
	ID   string
}
type availableTask = taskOption
type relatedTask taskOption

type TagWithTasks struct {
	tagUI
	RelatedTasks   []relatedTask
	AvailableTasks []availableTask
}

func newTagWithTasks(tag tagUI, relatedTasks []relatedTask, availableTasks []availableTask) TagWithTasks {
	return TagWithTasks{
		tagUI:          tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}
