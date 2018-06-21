package utils

import "os"

// Retrieve access tokens stored as environment variable
func GetEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}
