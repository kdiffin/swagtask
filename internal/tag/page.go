package tag

import "swagtask/internal/auth"

type tagsPage struct {
	TagsWithTasks []TagWithTasks
	Auth          auth.AuthenticatedPage
}

func newTagsPage(tagsWithTasks []TagWithTasks, authorized bool, pathToPfp string, username string) tagsPage {
	return tagsPage{
		TagsWithTasks: tagsWithTasks,
		Auth: auth.AuthenticatedPage{
			Authorized: authorized,
			UserUI: auth.UserUI{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		},
	}
}
