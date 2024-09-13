package universum

import (
	"reflect"
	"testing"
)

func TestConvertToStringBool(t *testing.T) {
	testCases := []struct {
		name      string
		input     map[string]interface{}
		expected  map[string]bool
		expectErr bool
	}{
		{
			name: "AllBoolValues",
			input: map[string]interface{}{
				"isEnabled": true,
				"isAdmin":   false,
			},
			expected: map[string]bool{
				"isEnabled": true,
				"isAdmin":   false,
			},
			expectErr: false,
		},
		{
			name: "NonBoolValueInMap",
			input: map[string]interface{}{
				"isEnabled": true,
				"age":       30,
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "EmptyMap",
			input:     map[string]interface{}{},
			expected:  map[string]bool{},
			expectErr: false,
		},
		{
			name: "SingleBoolValue",
			input: map[string]interface{}{
				"isValid": true,
			},
			expected: map[string]bool{
				"isValid": true,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := convertToStringBool(tc.input)

			if (err != nil) != tc.expectErr {
				t.Errorf("Test %s failed: expected error = %v, got %v", tc.name, tc.expectErr, err)
			}

			if !tc.expectErr && !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %s failed: expected %v, got %v", tc.name, tc.expected, actual)
			}
		})
	}
}
