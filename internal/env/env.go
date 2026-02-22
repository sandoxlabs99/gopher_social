package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok || val == "" {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)

	if !ok || val == "" {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)

	if err != nil {
		return fallback
	}

	return valAsInt
}

func GetDuration(key string, fallback string) time.Duration {
	val, ok := os.LookupEnv(key)

	if !ok || val == "" {
		duration, _ := time.ParseDuration(fallback)
		return duration
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		duration, _ := time.ParseDuration(fallback)
		return duration
	}

	return duration
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)

	if !ok || val == "" {
		return fallback
	}

	valAsBool, err := strconv.ParseBool(val)

	if err != nil {
		return fallback
	}

	return valAsBool
}
