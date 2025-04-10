package task_package

// task creation functions
type Task struct {
	Name string
	Idea string
	Id   int
    Tags []string
	Completed bool
}
type Tasks = []Task
type TasksContainer struct {
	Tasks Tasks
}
func NewTask(name string, idea string, id int, tags []string) Task {
	return Task{
		Name: name,
		Idea: idea,
		Id:   id,
        Tags: tags,
		Completed: false,
	}
}

func (p TasksContainer) HasIdea(idea string) bool {
    for _,task := range p.Tasks {
        if(task.Idea == idea) {
            return true
        }
    }

    return false
}

func NewHomePage() TasksContainer {
	return TasksContainer{
		Tasks: []Task{
			NewTask("ok bro", "read book", 2,[]string{"academia", "mental"}),
			NewTask("real", "do the dishes", 1, []string{"home", "physical"}),
		},
	}
}
