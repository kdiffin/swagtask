package task

import (
	"swagtask/internal/auth"

	"github.com/jackc/pgx/v5/pgtype"
)

type tasksPageFilters struct {
	SearchQuery string
	ActiveTag   string
}

func newTasksPageFilters(tagName string, taskName string) tasksPageFilters {

	return tasksPageFilters{
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
	taskWithTags
	Buttons taskPageButtons
	Auth    auth.AuthenticatedPage
}

func newTaskPage(task taskWithTags, prevButton, nextButton taskButton, authorized bool, pathToPfp string, username string) taskPage {
	return taskPage{
		taskWithTags: task,
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
type tasksPage struct {
	Tasks   []taskWithTags
	Filters tasksPageFilters
	Auth    auth.AuthenticatedPage
}

func newTasksPage(tasks []taskWithTags, filters tasksPageFilters,
	authorized bool, pathToPfp string, username string) tasksPage {
	return tasksPage{
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
