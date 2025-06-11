package envutils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvInt(key string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func GetEnvBool(key string, defaultVal bool) bool {
	if valStr, ok := os.LookupEnv(key); ok {
		valStr = strings.ToLower(valStr)
		return valStr == "true" || valStr == "1" || valStr == "yes"
	}
	return defaultVal
}

func GetEnvRune(key string, defaultVal rune) rune {
	if valStr, ok := os.LookupEnv(key); ok && len(valStr) > 0 {
		return rune(valStr[0])
	}
	return defaultVal
}
