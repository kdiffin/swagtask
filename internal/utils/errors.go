package utils

import (
	"errors"
	"net/http"
)

// ErrNoUpdateFields is returned when no fields are provided for an update (204 No Content)
var ErrNoUpdateFields = errors.New("no update fields provided for task")

// ErrNotFound is returned when a resource is not found (404 Not Found)
var ErrNotFound = errors.New("resource not found")

// ErrUnauthorized is returned when the user is not authenticated (401 Unauthorized)
var ErrUnauthorized = errors.New("user not authenticated")

// ErrForbidden is returned when the user lacks permission (403 Forbidden)
var ErrForbidden = errors.New("user lacks permission")

// ErrConflict is returned when there is a conflict with the current state (409 Conflict)
var ErrConflict = errors.New("conflict with current state")

// ErrBadRequest is returned for malformed requests or invalid input (400 Bad Request)
var ErrBadRequest = errors.New("bad request")

// ErrUnprocessable is returned for semantic errors in the request (422 Unprocessable Entity)
var ErrUnprocessable = errors.New("unprocessable entity")

var ErrInternalServer = errors.New("internal Server Error")

// knownErrors maps custom service errors to HTTP status codes.
// Add or update these as your service layer grows.
var knownErrors = map[error]int{
	ErrNotFound:       http.StatusNotFound,            // Resource not found (404 Not Found)
	ErrUnauthorized:   http.StatusUnauthorized,        // User not authenticated (401 Unauthorized)
	ErrForbidden:      http.StatusForbidden,           // User lacks permission (403 Forbidden)
	ErrConflict:       http.StatusConflict,            // Conflict with current state (409 Conflict)
	ErrBadRequest:     http.StatusBadRequest,          // Malformed request or invalid input (400 Bad Request)
	ErrUnprocessable:  http.StatusUnprocessableEntity, // Semantic errors in request (422 Unprocessable Entity)
	ErrInternalServer: http.StatusInternalServerError,
}
