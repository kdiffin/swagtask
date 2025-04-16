package pages

import "swagtask/database"

// tags.gohtml
type TagsPage struct {
	Tasks     []database.TaskWithTags
	Tag       string
	OtherTags []string
}

// task.gohtml
type TaskPage struct {
	Task    database.TaskWithTags
	Buttons struct {
		PrevButton TaskButton
		NextButton TaskButton
	}
}

type TaskButton struct {
	Id     int
	Name   string
	Exists bool
}

type IndexPage struct {
	Tasks []database.TaskWithTags
}
