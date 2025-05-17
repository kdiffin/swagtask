package tag

import "swagtask/internal/auth"

type tagsPage struct {
	TagsWithTasks []tagWithTasks
	Auth          auth.AuthenticatedPage
}

func newTagsPage(tagsWithTasks []tagWithTasks, authorized bool, pathToPfp string, username string) tagsPage {
	return tagsPage{
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
