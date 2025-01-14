package utils

import (
	"encoding/json"
	"fmt"
)

// GenerateCacheKey generates a cache key based on the input parameters
func GenerateCacheKey(prefix string, requiredParam any, optionalParams ...any) string {
	key := prefix
	key += fmt.Sprintf(":%v", requiredParam)
	for _, param := range optionalParams {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}

// Serialize marshals the input data into an array of bytes
func Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

// Deserialize unmarshals the input data into the output interface
func Deserialize(data []byte, output any) error {
	return json.Unmarshal(data, output)
}
