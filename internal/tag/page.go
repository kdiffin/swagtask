package tag

import "swagtask/internal/auth"

type TagsPage struct {
	TagsWithTasks []TagWithTasks
	Auth          auth.AuthenticatedPage
}

func NewTagsPage(tagsWithTasks []TagWithTasks, authorized bool, pathToPfp string, username string) TagsPage {
	return TagsPage{
		TagsWithTasks: tagsWithTasks,
		Auth: auth.AuthenticatedPage{
			Authorized: authorized,
			User: auth.UserUI{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		},
	}
}
