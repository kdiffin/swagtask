package service

import "errors"

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
