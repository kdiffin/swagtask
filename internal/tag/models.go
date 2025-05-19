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
type tagWithTasks struct {
	tagUI
	RelatedTasks   []relatedTask
	AvailableTasks []availableTask
}

func newTagWithTasks(tag tagUI, relatedTasks []relatedTask, availableTasks []availableTask) tagWithTasks {
	return tagWithTasks{
		tagUI:          tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}
