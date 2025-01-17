package env

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
	"time"
)

// GetString retrieves the value associated with the specified key from the .env file.
// If the key is not set, it panics with an error message.
func GetString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	return val
}

// GetInt retrieves the value associated with the specified key from the .env file,
// converting it to an integer. If the key is not set or the value cannot be parsed to an integer, it panics.
func GetInt(key string) int {
	val := GetString(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not an int", key))
	}
	return i
}

// GetDuration retrieves the value associated with the specified key from the .env file,
// converting it to a time.Duration. If the key is not set or the value cannot be parsed to a duration, it panics.
func GetDuration(key string) time.Duration {
	val := GetString(key)
	d, err := time.ParseDuration(val)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a duration", key))
	}
	return d
}
