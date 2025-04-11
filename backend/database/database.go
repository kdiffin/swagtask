package database

// task creation functions
type Task struct {
	Name      string
	Idea      string
	Id        int
	Tags      []string
	Completed bool
}
type Tasks = []Task
type Database struct {
	Tasks Tasks
}

func NewTask(name string, idea string, id int, tags []string) Task {
	return Task{
		Name:      name,
		Idea:      idea,
		Id:        id,
		Tags:      tags,
		Completed: false,
	}
}

func (p Database) HasIdea(idea string) bool {
	for _, task := range p.Tasks {
		if task.Idea == idea {
			return true
		}
	}

	return false
}

func NewDatabase() Database {
	return Database{
		Tasks: []Task{
			NewTask("ok bro", "read book", 2, []string{"academia", "mental"}),
			NewTask("real", "do the dishes", 1, []string{"home", "physical"}),
		},
	}
}
