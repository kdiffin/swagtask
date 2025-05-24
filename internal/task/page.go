package task

import (
	"swagtask/internal/auth"

	"github.com/jackc/pgx/v5/pgtype"
)

type TasksPageFilters struct {
	SearchQuery string
	ActiveTag   string
}

func newTasksPageFilters(tagName string, taskName string) TasksPageFilters {

	return TasksPageFilters{
		ActiveTag:   tagName,
		SearchQuery: taskName,
	}
}

// task individual page
type taskButton struct {
	ID     pgtype.UUID
	Name   string
	Exists bool
}
type taskPageButtons struct {
	PrevButton taskButton
	NextButton taskButton
}
type taskPage struct {
	TaskWithTags
	Buttons taskPageButtons
	Auth    auth.AuthenticatedPage
}

func newTaskPage(task TaskWithTags, prevButton, nextButton taskButton, authorized bool, pathToPfp string, username string) taskPage {
	return taskPage{
		TaskWithTags: task,
		Buttons: taskPageButtons{
			PrevButton: prevButton,
			NextButton: nextButton,
		},
		Auth: auth.AuthenticatedPage{
			Authorized: authorized,
			User: auth.UserUI{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		},
	}
}

// tasks.gohtml
type TasksPage struct {
	Tasks   []TaskWithTags
	Filters TasksPageFilters
	Auth    auth.AuthenticatedPage
}

func newTasksPage(tasks []TaskWithTags, filters TasksPageFilters,
	authorized bool, pathToPfp string, username string) TasksPage {
	return TasksPage{
		Tasks:   tasks,
		Filters: filters,
		Auth: auth.AuthenticatedPage{
			Authorized: authorized,
			User: auth.UserUI{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		},
	}

}
