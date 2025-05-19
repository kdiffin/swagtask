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
	task    taskWithTags
	Buttons taskPageButtons
	Auth    auth.AuthenticatedPage
}
