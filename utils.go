package universum

import (
	"fmt"
)

func convertToStringBool(input map[string]interface{}) (map[string]bool, error) {
	result := make(map[string]bool)

	for key, value := range input {
		boolVal, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("value for key '%s' is not of type bool", key)
		}
		result[key] = boolVal
	}

	return result, nil
}

func isWriteableDatatype(value interface{}) bool {
	switch value.(type) {
	case string,
		int, int8, int16, int32, uint32,
		uint, uint8, uint16, int64, uint64,
		float32, float64, bool:
		return true

	case []interface{}:
		return true

	default:
		return false
	}
}
