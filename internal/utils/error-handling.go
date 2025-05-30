package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
)

// add a message field here
func CheckError(w http.ResponseWriter, r *http.Request, err error) bool {
	for knownErr, status := range knownErrors {
		if errors.Is(err, knownErr) {
			if status == http.StatusUnauthorized {
				LogError("", err)
				// for htmx routes
				w.Header().Add("hx-redirect", "/login/")
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return true
			} else {
				LogError("", err)
				http.Error(w, http.StatusText(status), status)
				return true
			}
		}
	}
	if err != nil {
		LogError("", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}
	return false
}

func CheckErrorWebsocket(broadcast chan string, message string, err error) bool {
	for knownErr, status := range knownErrors {
		if errors.Is(err, knownErr) {
			if status == http.StatusUnauthorized {
				LogError(message, err)
				// for htmx routes
				broadcast <- fmt.Sprintf(`{
  						"type": "error",
  						"status": %v,
  						"message": "%v"
					}`, status, message)
				return true
			} else {
				LogError(message, err)
				broadcast <- fmt.Sprintf(`{
  						"type": "error",
  						"status": %v,
  						"message": "%v"
					}`, status, message)
				return true
			}
		}
	}
	if err != nil {
		LogError(message, err)
		broadcast <- fmt.Sprintf(`{
  						"type": "error",
  						"status": %v,
  						"message": "%v"
					}`, 0, message)
		return true
	}
	return false
}
func LogError(context string, err error) {
	if err == nil {
		return
	}
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown_file"
		line = 0
	}

	// Improve the log formatting for readability
	log.Printf(
		"----------------------------------\n"+
			"❌ [Error] %s:%d\n"+
			"Description: %s\n"+
			"Error: %+v\n"+
			"Stack Trace:\n%s"+
			"----------------------------------\n",
		file, line, context, err, debug.Stack(),
	)
}
