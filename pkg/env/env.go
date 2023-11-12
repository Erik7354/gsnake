package env

import (
	"os"
	"strconv"
)

func GetEnv(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}

func GetEnvInt(key string, def int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return parsed
}
