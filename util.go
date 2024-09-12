package universum

import (
	"fmt"
	"sort"
)

func sortMapKeys(m map[string]interface{}) map[string]interface{} {
	sorted := make(map[string]interface{})
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sorted[k] = m[k]
	}
	return sorted
}

func ConvertToStringBool(input map[string]interface{}) (map[string]bool, error) {
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
