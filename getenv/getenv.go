package getenv

import (
	"os"
	"strconv"
)

func String(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Bool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if value == "true" {
			return true
		}
	}

	return fallback
}

func Int(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		// is this a valid int?
		v, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return v
	}

	return fallback
}

func Int64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		// is this a valid int?
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return v
	}

	return fallback
}
