package tag

type tagUI struct {
	VaultID string
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
type TagWithTasks struct {
	Tag            tagUI
	RelatedTasks   []relatedTask
	AvailableTasks []availableTask
}

func newTagWithTasks(tag tagUI, relatedTasks []relatedTask, availableTasks []availableTask) TagWithTasks {
	return TagWithTasks{
		Tag:            tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}
