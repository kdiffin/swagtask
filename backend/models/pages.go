package models


type authenticatedPage struct {
	Authorized bool
	User
}
// tags.gohtml
type tagsPage struct {
	TagsWithTasks []TagWithTasks
	Auth authenticatedPage
}

func NewTagsPage(tagsWithTasks []TagWithTasks, authorized bool, pathToPfp string, username string) tagsPage {
	return tagsPage{
		TagsWithTasks: tagsWithTasks,
		Auth: authenticatedPage{
			Authorized: authorized,
			User: User{
				PathToPfp: pathToPfp,
				Username: username,
			},
		},
	
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
	Auth authenticatedPage
	
}

func NewTaskPage(task TaskWithTags, prevButton, nextButton TaskButton, authorized bool, pathToPfp string, username string) taskPage {
	return taskPage{
		Task: task,
		Buttons: TaskPageButtons{
			PrevButton: prevButton,
			NextButton: nextButton,
		},
		Auth: authenticatedPage{
			Authorized: authorized,
			User: User{
				PathToPfp: pathToPfp,
				Username: username,
			},
		},
	
	}
}

// tasks.gohtml
type TasksPage struct {
	Tasks []TaskWithTags
	Filters *TasksPageFilters
	Auth authenticatedPage
}

func NewTasksPage(tasks []TaskWithTags,tagFilter *string, taskFilter *string, authorized bool, pathToPfp string, username string) TasksPage {
	var filters *TasksPageFilters 
	if tagFilter != nil && taskFilter != nil {
		filtersVal := NewTasksPageFilters(*tagFilter, *taskFilter)
		filters = &filtersVal		
	} 
	return TasksPage{
		Tasks: tasks,
		Filters: filters,
		Auth: authenticatedPage{
			Authorized: authorized,
			User: User{
				PathToPfp: pathToPfp,
				Username: username,
			},
		},
	
	}
	
}
