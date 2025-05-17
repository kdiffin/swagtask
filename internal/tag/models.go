package tag

import "swagtask/internal/auth"

type tagAuthor = auth.Author

type tagUI struct {
	VaultID string
	Author  tagAuthor
	Name    string
	ID      string
}
type availableTask struct {
	Name string
	ID   string
}
type relatedTask struct {
	Name string
	ID   string
}
type tagWithTasks struct {
	Tag            tagUI
	RelatedTasks   []relatedTask
	AvailableTasks []availableTask
}

func newTagWithTasks(tag tagUI, relatedTasks []relatedTask, availableTasks []availableTask) tagWithTasks {
	return tagWithTasks{
		Tag:            tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}
