package task

type TasksPageFilters struct {
	SearchQuery string
	ActiveTag   string
}

func NewTasksPageFilters(tagName string, taskName string) TasksPageFilters {
	return TasksPageFilters{
		ActiveTag:   tagName,
		SearchQuery: taskName,
	}
}
