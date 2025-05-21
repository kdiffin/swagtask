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

type vaultUI struct {
	ID          string
	Name        string
	Description string
	Author      auth.Author
	Locked      bool
	CreatedAt   string
	Kind        kind
}

func newUIvault(vault db.Vault, author auth.Author) vaultUI {

	return vaultUI{
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
}

// type availableCollaborator = collaboratorOption
type relatedCollaborator collaboratorOption
type vaultWithCollaborators struct {
	vaultUI
	RelatedCollaborators []relatedCollaborator
	// AvailableCollaborators []availableCollaborator
}

func newVaultWithCollaborators(vault vaultUI, relatedCollaborators []relatedCollaborator) vaultWithCollaborators {
	return vaultWithCollaborators{
		vaultUI:              vault,
		RelatedCollaborators: relatedCollaborators,
	}
}
