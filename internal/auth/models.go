package auth

// user
type UserUI struct {
	ID        string
	PathToPfp string
	Username  string
}

type AuthenticatedPage struct {
	Authorized bool
	User       UserUI
}
