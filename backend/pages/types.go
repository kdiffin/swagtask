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
	Task           database.TaskWithTags
	PrevTaskExists bool
	NextTaskExists bool
	PrevId         int
	NextId         int
}

type IndexPage struct {
	Tasks []database.TaskWithTags
}
