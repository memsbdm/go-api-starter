package utils

import (
	"encoding/json"
	"fmt"
)

// GenerateCacheKey generates a cache key using a prefix and input parameters.
func GenerateCacheKey(prefix string, requiredParam any, optionalParams ...any) string {
	key := prefix
	key += fmt.Sprintf(":%v", requiredParam)
	for _, param := range optionalParams {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}

// Serialize marshals the input data into a byte array.
// It converts the provided data structure into JSON format.
func Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

// Deserialize unmarshal the input byte array into the specified output interface.
// It converts JSON bytes back into the original data structure.
func Deserialize(data []byte, output any) error {
	return json.Unmarshal(data, output)
}
