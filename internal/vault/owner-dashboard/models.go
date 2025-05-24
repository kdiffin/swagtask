package vault

import (
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"
)

type kind string

const (
	kindDefault       kind = "default"
	kindCollaborative kind = "collaborative"
)

type VaultUI struct {
	ID          string
	Name        string
	Description string
	Author      auth.Author
	Locked      bool
	CreatedAt   string
	Kind        kind
}

func newUIvault(vault db.Vault, author auth.Author) VaultUI {

	return VaultUI{
		ID:          vault.ID.String(),
		Author:      author,
		Name:        vault.Name,
		Description: vault.Description,
		Locked:      vault.Locked,
		CreatedAt:   utils.BrowserFormattedtTime(vault.CreatedAt),
		Kind:        kind(vault.Kind),
	}
}

// ---- FOR UI ----
// vaults
type collaboratorOption struct {
	Name      string
	PathToPfp string
	Role      string
}

// type availableCollaborator = collaboratorOption
type relatedCollaborator collaboratorOption
type vaultWithCollaborators struct {
	VaultUI
	RelatedCollaborators []relatedCollaborator
	// AvailableCollaborators []availableCollaborator
}

func newVaultWithCollaborators(vault VaultUI, relatedCollaborators []relatedCollaborator) vaultWithCollaborators {
	return vaultWithCollaborators{
		VaultUI:              vault,
		RelatedCollaborators: relatedCollaborators,
	}
}

// for collaborators
type role string

const (
	roleViewer       kind = "viewer"
	roleCollaborator kind = "collaborator"
	roleOwner        kind = "owner"
)
