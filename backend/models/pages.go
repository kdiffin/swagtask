package models

// tags.gohtml
type tagsPage struct {
	TagsWithTasks []TagWithTasks
}

func NewTagsPage(allTags []TagWithTasks) tagsPage {
	return tagsPage{
		TagsWithTasks: allTags,
	}
}

// task.gohtml
type TaskButton struct {
	Id     int
	Name   string
	Exists bool
}
type TaskPageButtons struct {
	PrevButton TaskButton
	NextButton TaskButton
}
type taskPage struct {
	Task    TaskWithTags
	Buttons TaskPageButtons
}

func NewTaskPage(task TaskWithTags, prevButton, nextButton TaskButton) taskPage {
	return taskPage{
		Task: task,
		Buttons: TaskPageButtons{
			PrevButton: prevButton,
			NextButton: nextButton,
		},
	}
}

// tasks.gohtml
type TasksPage struct {
	Tasks []TaskWithTags
}

func NewTasksPage(tasks []TaskWithTags) TasksPage {
	return TasksPage{
		Tasks: tasks,
	}
}
