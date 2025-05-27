package vault

import (
	"swagtask/internal/auth"
	"swagtask/internal/tag"
	"swagtask/internal/task"
)

// collaborative

type CollaboratorUI = collaboratorOption
type UserVaultUI struct {
	PathToPfp string
	Username  string
	Role      string
}

type vaultTasksPage struct {
	User          UserVaultUI
	Vault         VaultUI
	Collaborators []CollaboratorUI

	task.TasksPage
}

func NewVaultTasksPage(
	tasks []task.TaskWithTags,
	filters task.TasksPageFilters,
	authorized bool,
	pathToPfp string,
	vault VaultUI,
	User UserVaultUI,
	collaborators []CollaboratorUI,
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
		Vault:         vault,
		User:          User,
		Collaborators: collaborators,
	}
}

type vaultTagsPage struct {
	tag.TagsPage
	User          UserVaultUI
	Collaborators []CollaboratorUI
}

func newVaultTagsPage(
	tags []tag.TagWithTasks,
	authorized bool,
	pathToPfp string,
	User UserVaultUI,
	collaborators []CollaboratorUI,
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
	User          UserVaultUI
	Vault         VaultUI
	Collaborators []CollaboratorUI
	Auth          auth.AuthenticatedPage
}

func NewVaultPage(
	authorized bool,
	User UserVaultUI,
	pathToPfp string,
	username string,
	vault VaultUI,
	collaborators []CollaboratorUI,
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
