package vault

import "swagtask/internal/auth"

type vaultsPage struct {
	Vaults []VaultWithCollaborators
	Auth   auth.AuthenticatedPage
}

func NewVaultsPage(vaults []VaultWithCollaborators,
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
