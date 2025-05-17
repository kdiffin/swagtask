package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	db "swagtask/internal/db/generated"
)

// TODO: write this well
func HandleSignup(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hash, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	defaultVaultId, errDefaultVault := queries.CreateDefaultVault(r.Context(), db.CreateDefaultVaultParams{
		Name:        "Default",
		Description: "This is your default vault. Only you can access this.",
	})
	if errDefaultVault != nil {
		http.Error(w, "Error creating default vault"+http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = queries.CreateUser(r.Context(), db.CreateUserParams{
		Username:       username,
		PasswordHash:   hash,
		DefaultVaultID: defaultVaultId,
	})
	if err != nil {
		http.Error(w, "User exists", 400)
		return
	}

	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
func HandleLogin(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	credentials, err := queries.GetUserCredentials(r.Context(), username)
	if err != nil || !CheckPasswordHash(password, credentials.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := generateSessionID()
	errSesh := queries.CreateSession(r.Context(), db.CreateSessionParams{ID: sessionID, UserID: credentials.ID})
	if errSesh != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/tasks/", http.StatusSeeOther)
}

func HandleLogout(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		errDeleteCookie := queries.DeleteSession(r.Context(), cookie.Value)
		if errDeleteCookie != nil {
			http.Error(w, "Error deleting cookie, try logging out again", http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Secure: true,
		MaxAge: -1,
		Path:   "/",
	})
	http.Redirect(w, r, "/login/", http.StatusSeeOther)
}
