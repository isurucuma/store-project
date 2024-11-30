package config

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

func GetApplicationPort() int {
	port, err := strconv.Atoi(GetEnv("APPLICATION_PORT", "8080"))
	if err != nil {
		return 8080
	}
	return port
}
