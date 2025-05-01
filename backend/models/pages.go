package models

// tags.gohtml
type tagsPage struct {
	TagsWithTasks []TagWithTasks	
}

func NewTagsPage(tagsWithTasks []TagWithTasks) tagsPage {
	return tagsPage{
		TagsWithTasks: tagsWithTasks,
	
	}
}

// task.gohtml
type TaskButton struct {
	ID    int32
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
	Filters *TasksPageFilters
}

func NewTasksPage(tasks []TaskWithTags,tagFilter *string, taskFilter *string) TasksPage {
	var filters *TasksPageFilters 
	if tagFilter != nil && taskFilter != nil {
		filtersVal := NewTasksPageFilters(*tagFilter, *taskFilter)
		filters = &filtersVal		
	} 
	return TasksPage{
		Tasks: tasks,
		Filters: filters,
	}
}
