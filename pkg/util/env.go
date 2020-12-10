package util

import (
	"os"
	"strconv"
)

func EnvString(key string, def string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return env
}

func EnvBool(key string, def bool) bool {
	env, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	val, err := strconv.ParseBool(env)
	if err != nil {
		return def
	}

	return val
}
