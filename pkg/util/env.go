package util

import (
	"os"
	"strconv"
)

func EnvString(key string, val string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return val
	}

	return val
}

func EnvBool(key string, val bool) bool {
	rawVal, ok := os.LookupEnv(key)
	if !ok {
		return val
	}
	val, err := strconv.ParseBool(rawVal)
	if err != nil {
		return val
	}

	return val
}
