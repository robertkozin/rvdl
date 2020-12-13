package util

import (
	"os"
	"strconv"
)

func EnvString(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}

	return val
}

func EnvBool(key string, defaultVal bool) bool {
	rawVal, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	val, err := strconv.ParseBool(rawVal)
	if err != nil {
		return defaultVal
	}

	return val
}
