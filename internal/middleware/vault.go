package middleware

import (
	"context"
	"log"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"
)

const vaultContextKey = contextKey("vault")

type VaultInfo struct {
	ID string
}

func VaultIDFromContext(ctx context.Context) (string, error) {
	vaultInfo, ok := ctx.Value(vaultContextKey).(*VaultInfo)
	if !ok {
		log.Println("vault ID not found in context")
		return "", utils.ErrInternalServer
	}
	return vaultInfo.ID, nil
}

// Middleware to extract vault ID from the URL path
func HandlerWithVaultIdFromPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vaultID := r.PathValue("vaultId")
		if vaultID == "" {
			utils.CheckError(w, r, utils.ErrBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), vaultContextKey, &VaultInfo{ID: vaultID})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware to extract vault ID from the user's default
func HandlerWithVaultIdFromUser(queries *db.Queries, next http.Handler) http.Handler {
	return HandlerWithUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if utils.CheckError(w, r, err) {
			log.Println("error getting user in vaults middleware")
			return
		}

		ctx := context.WithValue(r.Context(), vaultContextKey, &VaultInfo{ID: user.DefaultVaultID})
		next.ServeHTTP(w, r.WithContext(ctx))

	}))
}
