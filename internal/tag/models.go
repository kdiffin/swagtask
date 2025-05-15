package tag

import (
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
)

type availableTask struct {
	Name string
	ID   int32
}
type relatedTask struct {
	Name string
	ID   int32
}
type tagWithTasks struct {
	db.Tag
	RelatedTasks   []relatedTask
	AvailableTasks []availableTask
}

func newTagWithTasks(tag db.Tag, relatedTasks []relatedTask, availableTasks []availableTask) TagWithTasks {
	return tagWithTasks{
		Tag:            tag,
		RelatedTasks:   relatedTasks,
		AvailableTasks: availableTasks,
	}
}

type tagsPage struct {
	TagsWithTasks []TagWithTasks
	Auth          auth.AuthenticatedPage
}

func newTagsPage(tagsWithTasks []TagWithTasks, authorized bool, pathToPfp string, username string) tagsPage {
	return tagsPage{
		TagsWithTasks: tagsWithTasks,
		Auth: auth.AuthenticatedPage{
			Authorized: authorized,
			User: auth.User{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		},
	}
}
