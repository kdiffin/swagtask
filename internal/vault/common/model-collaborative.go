package vault

import (
	"swagtask/internal/tag"
	"swagtask/internal/task"
)

// collaborative

type CollaboratorUI = collaboratorOption
type UserVaultUI struct {
	PathToPfp  string
	Username   string
	Role       string
	Authorized bool
}

type vaultTasksPage struct {
	Auth          UserVaultUI
	Vault         VaultUI
	Collaborators []CollaboratorUI

	task.TasksPage
}

func NewVaultTasksPage(
	tasks []task.TaskWithTags,
	filters task.TasksPageFilters,

	vault VaultUI,
	User UserVaultUI,
	collaborators []CollaboratorUI,
) vaultTasksPage {

	return vaultTasksPage{
		TasksPage: task.TasksPage{
			Tasks:   tasks,
			Filters: filters,
		},
		Vault:         vault,
		Auth:          User,
		Collaborators: collaborators,
	}
}

type vaultTaskPage struct {
	Auth          UserVaultUI
	Vault         VaultUI
	Collaborators []CollaboratorUI

	task.TaskPage
}

func NewVaultTaskPage(
	taskWithTags task.TaskWithTags,

	buttons task.TaskPageButtons,

	vault VaultUI,
	User UserVaultUI,
	collaborators []CollaboratorUI,
) vaultTaskPage {

	return vaultTaskPage{
		TaskPage: task.TaskPage{
			TaskWithTags: taskWithTags,
			Buttons:      buttons,
		},
		Vault:         vault,
		Auth:          User,
		Collaborators: collaborators,
	}
}

type vaultTagsPage struct {
	tag.TagsPage
	Vault VaultUI

	Auth          UserVaultUI
	Collaborators []CollaboratorUI
}

func NewVaultTagsPage(
	tags []tag.TagWithTasks,
	Vault VaultUI,
	User UserVaultUI,
	collaborators []CollaboratorUI,
) vaultTagsPage {

	return vaultTagsPage{
		TagsPage: tag.TagsPage{
			TagsWithTasks: tags,
		},
		Vault:         Vault,
		Auth:          User,
		Collaborators: collaborators,
	}
}

type vaultHomePage struct {
	Auth          UserVaultUI
	Vault         VaultUI
	Collaborators []CollaboratorUI
}

func NewVaultPage(

	User UserVaultUI,
	vault VaultUI,
	collaborators []CollaboratorUI,
) vaultHomePage {

	return vaultHomePage{
		Collaborators: collaborators,
		Auth:          User,

		Vault: vault,
	}
}
