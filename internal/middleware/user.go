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
func UserFromContext(ctx context.Context) (*db.User, bool) {
	user, ok := ctx.Value(userContextKey).(*db.User)
	return user, ok
}

func getUserIDFromRequest(queries *db.Queries, r *http.Request) (pgtype.UUID, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("error getting cookie: %w", err)
	}

	sesh, errSesh := queries.GetSessionValues(r.Context(), cookie.Value)
	if errSesh != nil {

		return pgtype.UUID{}, errSesh
	}
	return sesh.UserID, nil
}

func getUserInfoFromSessionId(queries *db.Queries, r *http.Request) (*auth.UserUI, error) {
	userId, err := getUserIDFromRequest(queries, r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrUnauthorized, err)
	}

	userDb, errUser := queries.GetUserInfo(r.Context(), userId)
	if errUser != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrUnauthorized, errUser)
	}
	user := auth.UserUI{
		ID:        userId.String(),
		PathToPfp: userDb.PathToPfp.String,
		Username:  userDb.Username,
	}
	return &user, nil

}

func WithUser(queries *db.Queries, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserInfoFromSessionId(queries, r)
		if err != nil {
			http.Redirect(w, r, "/login/", http.StatusUnauthorized)
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
