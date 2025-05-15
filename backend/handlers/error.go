package handlers

import (
	"errors"
	"net/http"
	"swagtask/service"
	"swagtask/utils"
)

// knownErrors maps custom service errors to HTTP status codes.
// Add or update these as your service layer grows.
var knownErrors = map[error]int{
	service.ErrNotFound:      http.StatusNotFound,            // Resource not found (404 Not Found)
	service.ErrUnauthorized:  http.StatusUnauthorized,        // User not authenticated (401 Unauthorized)
	service.ErrForbidden:     http.StatusForbidden,           // User lacks permission (403 Forbidden)
	service.ErrConflict:      http.StatusConflict,            // Conflict with current state (409 Conflict)
	service.ErrBadRequest:    http.StatusBadRequest,          // Malformed request or invalid input (400 Bad Request)
	service.ErrUnprocessable: http.StatusUnprocessableEntity, // Semantic errors in request (422 Unprocessable Entity)
}

func checkErrors(w http.ResponseWriter, r *http.Request, err error) bool {
	for knownErr, status := range knownErrors {
		if errors.Is(err, knownErr) {
			if status == http.StatusUnauthorized {
				utils.LogError("", err)
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return true
			} else {
				utils.LogError("", err)
				http.Error(w, http.StatusText(status), status)
				return true
			}
		}
	}
	if err != nil {
		utils.LogError("", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}
	return false
}
