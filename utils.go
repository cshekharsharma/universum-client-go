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
