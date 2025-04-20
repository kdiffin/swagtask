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
type TagRelationOption struct {
	Name string
	Id   int
}
type TagWithTasks struct {
	Tag
	AvailableTasks []TagRelationOption
}

func NewTagWithTasks(id int, name string, availableTasks []TagRelationOption) TagWithTasks {
	return TagWithTasks{
		Tag: Tag{
			Id:   id,
			Name: name,
		},
		AvailableTasks: availableTasks,
	}
}
