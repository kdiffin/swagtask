package vault

import (
	"swagtask/internal/auth"
	"swagtask/internal/tag"
	"swagtask/internal/task"
	vault "swagtask/internal/vault/owner-dashboard"
)

type userVaultUI struct {
	PathToPfp string
	Username  string
	Role      string
}

type vaultTasksPage struct {
	User          userVaultUI
	Collaborators []collaboratorUI

	task.TasksPage
}

func newVaultTasksPage(
	tasks []task.TaskWithTags,
	filters task.TasksPageFilters,
	authorized bool,
	pathToPfp string,
	User userVaultUI,
	collaborators []collaboratorUI,
	username string) vaultTasksPage {

	return vaultTasksPage{
		TasksPage: task.TasksPage{
			Tasks:   tasks,
			Filters: filters,
			Auth: auth.AuthenticatedPage{
				Authorized: authorized,
				User: auth.UserUI{
					PathToPfp: pathToPfp,
					Username:  username,
				},
			},
		},
		User:          User,
		Collaborators: collaborators,
	}
}

type vaultTagsPage struct {
	tag.TagsPage
	User          userVaultUI
	Collaborators []collaboratorUI
}

func newVaultTagsPage(
	tags []tag.TagWithTasks,
	authorized bool,
	pathToPfp string,
	User userVaultUI,
	collaborators []collaboratorUI,
	username string) vaultTagsPage {

	return vaultTagsPage{
		TagsPage: tag.TagsPage{
			TagsWithTasks: tags,
			Auth: auth.AuthenticatedPage{
				Authorized: authorized,
				User: auth.UserUI{
					PathToPfp: pathToPfp,
					Username:  username,
				},
			},
		},
		User:          User,
		Collaborators: collaborators,
	}
}

type vaultHomePage struct {
	User          userVaultUI
	Vault         vault.VaultUI
	Collaborators []collaboratorUI
	Auth          auth.AuthenticatedPage
}

func newVaultPage(
	authorized bool,
	User userVaultUI,
	pathToPfp string,
	username string,
	vault vault.VaultUI,
	collaborators []collaboratorUI,
) vaultHomePage {

	return vaultHomePage{
		Collaborators: collaborators,
		User:          User,
		Auth: auth.AuthenticatedPage{
			Authorized: authorized,
			User: auth.UserUI{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		},
		Vault: vault,
	}
}
