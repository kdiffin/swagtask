package utils

import (
	"fmt"
	"log"
	"runtime"
)

type FallbackString struct {
}

func StringWithFallback(main string, fallback string) string {
	if len(main) > 0 {
		return main
	} else {
		return fallback
	}
}

func PrintList[T any](list []T) {
	for i, v := range list {
		fmt.Printf("[%d]: %v\n", i, v)
	}
}



func LogError(context string, err error) {
	if err == nil {
		return
	}
	_, file, line, ok := runtime.Caller(1)
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
