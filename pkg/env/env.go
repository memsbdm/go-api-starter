package env

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

// GetString takes a key and returns the associated value in .env file. If there is no key it panics.
func GetString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	return val
}

// GetInt takes a key and returns the associated value in .env file converted in an integer.
// If there is no key or the value cannot be parsed to an integer it panics.
func GetInt(key string) int {
	val := GetString(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not an int", key))
	}
	return i
}
