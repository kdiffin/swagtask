package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"swagtask/auth"
	db "swagtask/db/generated"
)

// TODO: write this well
func handleSignup(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

    hash, err := auth.HashPassword(password)
    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }

    err = queries.CreateUser(r.Context(), db.CreateUserParams{
		Username: username,
		PasswordHash: hash,
	})
    if err != nil {
        http.Error(w, "User exists", 400)
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}



var sessions = map[string]int{} // sessionID -> userID (in-memory)

func generateSessionID() string {
    b := make([]byte, 32)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func handleLogin(queries *db.Queries, w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

	credentials, err := queries.GetUserCredentials(r.Context(), username)
	if err != nil || !auth.CheckPasswordHash(password, credentials.PasswordHash) {
        http.Error(w, "Invalid credentials", 401)
        return
    }

    sessionID := generateSessionID()
    sessions[sessionID] = int(credentials.ID)

    cookie := http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Secure: true,
        HttpOnly: true,
        Path:     "/",
    }
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}


func handleLogout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == nil {
        delete(sessions, cookie.Value)
    }

    http.SetCookie(w, &http.Cookie{
        Name:   "session_id",
        Secure: true,
        MaxAge: -1,
        Path:   "/",
    })
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}