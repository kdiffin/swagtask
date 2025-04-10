package utils

type FallbackString struct {
}

func StringWithFallback(main string, fallback string) string {
	if len(main) > 0 {
		return main
	} else {
		return fallback
	}
}
