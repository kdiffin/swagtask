package utils

import "fmt"

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
