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

// composed types
type TaskWithTags struct {
	Task
	Tags []Tag
}
