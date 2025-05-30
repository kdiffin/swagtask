package tag

import "swagtask/internal/auth"

type tagAuthor = auth.Author

type TagUI struct {
	VaultID   string
	Author    tagAuthor
	Name      string
	CreatedAt string

	ID string
}
type taskOption struct {
	Name string
	ID   string
}
type availableTask = taskOption
type relatedTask taskOption

type TagWithTasks struct {
	TagUI
	RelatedTasks   []relatedTask
	AvailableTasks []availableTask
}

func newTagWithTasks(tag TagUI, relatedTasks []relatedTask, availableTasks []availableTask) TagWithTasks {
	return TagWithTasks{
		TagUI:          tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}
