package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
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

// GetOptionalString retrieves the value associated with the specified key from the .env file.
// If the key is not set, it returns an empty string.
func GetOptionalString(key string) string {
	return os.Getenv(key)
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

// GetOptionalInt retrieves the value associated with the specified key from the .env file,
// converting it to an integer. If the key is not set or the value cannot be parsed to an integer, it returns 0.
func GetOptionalInt(key string) int {
	val := GetOptionalString(key)
	if val == "" {
		return 0
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
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

// GetOptionalDuration retrieves the value associated with the specified key from the .env file,
// converting it to a time.Duration. If the key is not set or the value cannot be parsed to a duration, it returns 0.
func GetOptionalDuration(key string) time.Duration {
	val := GetOptionalString(key)
	if val == "" {
		return 0
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0
	}
	return d
}

// GetFloat64 retrieves the value associated with the specified key from the .env file,
// converting it to a float64. If the key is not set or the value cannot be parsed to a float64, it panics.
func GetFloat64(key string) float64 {
	val := GetString(key)
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a float64", key))
	}
	return f
}

// GetOptionalFloat64 retrieves the value associated with the specified key from the .env file,
// converting it to a float64. If the key is not set or the value cannot be parsed to a float64, it returns 0.
func GetOptionalFloat64(key string) float64 {
	val := GetOptionalString(key)
	if val == "" {
		return 0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return f
}
