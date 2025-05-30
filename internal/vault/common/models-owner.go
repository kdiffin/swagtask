package vault

import (
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"
)

// owner
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

func NewUIvault(vault db.Vault, author auth.Author) VaultUI {

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
	Active    bool
}

// type availableCollaborator = collaboratorOption
type RelatedCollaborator collaboratorOption
type VaultWithCollaborators struct {
	VaultUI
	RelatedCollaborators []RelatedCollaborator
	// AvailableCollaborators []availableCollaborator
}

func newVaultWithCollaborators(vault VaultUI, relatedCollaborators []RelatedCollaborator) VaultWithCollaborators {
	return VaultWithCollaborators{
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
