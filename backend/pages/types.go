package pages

import "swagtask/database"

// tags.gohtml
type tagsPage struct {
	TagsWithTasks []database.TagWithTasks
}

func NewTagsPage(allTags []database.TagWithTasks) tagsPage {
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
	Task    database.TaskWithTags
	Buttons TaskPageButtons
}

func NewTaskPage(task database.TaskWithTags, prevButton, nextButton TaskButton) taskPage {
	return taskPage{
		Task: task,
		Buttons: TaskPageButtons{
			PrevButton: prevButton,
			NextButton: nextButton,
		},
	}
}

// tasks.gohtml
type tasksPage struct {
	Tasks []database.TaskWithTags
}

func NewTasksPage(tasks []database.TaskWithTags) tasksPage {
	return tasksPage{
		Tasks: tasks,
	}
}
