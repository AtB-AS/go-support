package config

import (
	"fmt"
	"os"
)

func GetRequiredEnv(key string) string {
	value, set := os.LookupEnv(key)
	if !set {
		panic(fmt.Sprintf("Environment variable '%s' is not set", key))
	}
	if value == "" {
		panic(fmt.Sprintf("No value set for environment variable '%s'", key))
	}

	return value
}
