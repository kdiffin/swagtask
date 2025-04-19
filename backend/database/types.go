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
	Tags    []Tag
	AllTags []Tag
}

func NewTaskWithTags(task Task, tags []Tag, allTags []Tag) TaskWithTags {
	return TaskWithTags{
		Task:    task,
		Tags:    tags,
		AllTags: allTags,
	}
}

// tags
type TagRelationOption struct {
	Name string
	Id   int
}
type TagWithTasks struct {
	Tag
	AllTasks []TagRelationOption
}

func NewTagWithTasks(id int, name string, allTasks []TagRelationOption) TagWithTasks {
	return TagWithTasks{
		Tag: Tag{
			Id:   id,
			Name: name,
		},
		AllTasks: allTasks,
	}
}
