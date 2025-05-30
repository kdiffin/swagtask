package auth

// user
type User struct {
	ID             string
	PathToPfp      string
	Username       string
	DefaultVaultID string
}

type Author struct {
	PathToPfp string
	Username  string
}

type UserUI struct {
	PathToPfp string
	Username  string
}

type AuthenticatedPage struct {
	Authorized bool
	User       UserUI
}
