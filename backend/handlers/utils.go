package handlers

import "net/http"

func getUserIDFromRequest(r *http.Request) (int, bool) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        return 0, false
    }

    userID, ok := sessions[cookie.Value]
    return userID, ok
}