package utils

import "os"

// Getenv returns env var or default val
func Getenv(envName string, defaultVal string) string {
	if u := os.Getenv(envName); u != "" {
		return u
	}
	return defaultVal
}
