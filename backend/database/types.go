package database

// db types
type Task struct {
	Name      string
	Idea      string
	Id        int
	Completed bool
}
type Tag struct {
	Id   int
	Name string
}

// ---- FOR UI ----
// tasks
type TaskWithTags struct {
	Task
	Tags          []Tag
	AvailableTags []Tag
}

func NewTaskWithTags(task Task, tags []Tag, availableTags []Tag) TaskWithTags {
	return TaskWithTags{
		Task:          task,
		Tags:          tags,
		AvailableTags: availableTags,
	}
}

// tags
type AvailableTask struct {
	Name string
	Id   int
}
type RelatedTask struct {
	Name string
	Id   int
}
type TagWithTasks struct {
	Tag
	RelatedTasks   []RelatedTask
	AvailableTasks []AvailableTask
}

func NewTagWithTasks(tag Tag, relatedTasks []RelatedTask, availableTasks []AvailableTask) TagWithTasks {
	return TagWithTasks{
		Tag:            tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}
