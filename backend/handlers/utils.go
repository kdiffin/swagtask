package handlers

import (
	"fmt"
	"net/http"
	db "swagtask/db/generated"
)

func getUserIDFromRequest(queries *db.Queries,r *http.Request) (int, error) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        return 0, fmt.Errorf("error getting cookie: %w", err)
    }

    sesh, errSesh := queries.GetSessionValues(r.Context(), cookie.Value)
    if errSesh != nil {
        return 0, errSesh
    }
    return int(sesh.UserID.Int32), nil
}