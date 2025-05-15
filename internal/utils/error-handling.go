package utils

import (
	"errors"
	"log"
	"net/http"
	"runtime"
)

func CheckError(w http.ResponseWriter, r *http.Request, err error) bool {
	for knownErr, status := range knownErrors {
		if errors.Is(err, knownErr) {
			if status == http.StatusUnauthorized {
				LogError("", err)
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
			"Error: %v\n"+
			"----------------------------------\n",
		file, line, context, err,
	)
}
