package middleware

import (
	"context"
	"fmt"
	"net/http"
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey string

const userContextKey = contextKey("user")

// Helper to get user from context
func UserFromContext(ctx context.Context) (*auth.UserUI, error) {
	user, ok := ctx.Value(userContextKey).(*auth.UserUI)
	if !ok {
		return nil, utils.ErrUnauthorized
	}
	return user, nil
}

func getUserIDFromRequest(queries *db.Queries, r *http.Request) (pgtype.UUID, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.LogError("error at COOKIE:", err)
		return pgtype.UUID{}, fmt.Errorf("error getting cookie: %w", err)
	}

	sesh, errSesh := queries.GetSessionValues(r.Context(), cookie.Value)
	if errSesh != nil {
		utils.LogError("error at session value from cookie:", err)
		return pgtype.UUID{}, errSesh
	}
	return sesh.UserID, nil
}

func getUserInfoFromSessionId(queries *db.Queries, r *http.Request) (*auth.UserUI, error) {
	userId, err := getUserIDFromRequest(queries, r)
	if err != nil {
		utils.LogError("error at userid from request:", err)
		return nil, fmt.Errorf("%w: %v", utils.ErrUnauthorized, err)
	}

	userDb, errUser := queries.GetUserInfo(r.Context(), userId)
	if errUser != nil {
		utils.LogError("error at user info:", errUser)
		return nil, fmt.Errorf("%w: %v", utils.ErrUnauthorized, errUser)
	}
	user := auth.UserUI{
		ID:        userId.String(),
		PathToPfp: userDb.PathToPfp.String,
		Username:  userDb.Username,
	}
	return &user, nil

}

func HandlerWithUser(queries *db.Queries, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserInfoFromSessionId(queries, r)
		if err != nil {
			utils.LogError("error at session value from cookie:", err)
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
