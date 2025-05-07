package handlers

import (
	"fmt"
	"net/http"
	db "swagtask/db/generated"
	"swagtask/models"
	"swagtask/service"
)

func getUserIDFromRequest(queries *db.Queries,r *http.Request) (int32, error) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        return 0, fmt.Errorf("error getting cookie: %w", err)
    }

    sesh, errSesh := queries.GetSessionValues(r.Context(), cookie.Value)
    if errSesh != nil {
        return 0, errSesh
    }
    return sesh.UserID.Int32, nil
}

func getUserInfoFromSessionId(queries *db.Queries,r *http.Request) (*models.User, error) {
    userId, err := getUserIDFromRequest(queries, r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", service.ErrUnauthorized, err)
	}

	fmt.Println("USER ID:", userId)
    userDb, errUser := queries.GetUserInfo(r.Context(), userId)
    if errUser != nil {
        return nil, fmt.Errorf("%w: %v", service.ErrUnauthorized, errUser)
    }
    user := models.User{
        PathToPfp: userDb.PathToPfp.String,
        Username: userDb.Username,
    }
    return &user, nil

}