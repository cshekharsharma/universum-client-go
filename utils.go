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
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true

	case []string, []bool,
		[]int, []int8, []int16, []int32, []int64,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64, []interface{}:
		return true

	default:
		return false
	}
}
