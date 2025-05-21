package vault

import (
	"swagtask/internal/auth"
)

// tasks.gohtml
type vaultsPage struct {
	Vaults []vaultWithCollaborators
	Auth   auth.AuthenticatedPage
}

func newVaultsPage(vaults []vaultWithCollaborators,
	authorized bool, pathToPfp string, username string) vaultsPage {
	return vaultsPage{
		Vaults: vaults,
		Auth: auth.AuthenticatedPage{
			Authorized: authorized,
			User: auth.UserUI{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		},
	}

}
