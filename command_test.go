package universum

import (
	"errors"
	"testing"
)

func TestToCommandResult(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := []interface{}{"some-value", int64(200), "OK"}
		result, err := toCommandResult(input)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.value != "some-value" {
			t.Errorf("Expected value to be 'some-value', got %v", result.value)
		}
		if result.code != 200 {
			t.Errorf("Expected code to be 200, got %v", result.code)
		}
		if result.message != "OK" {
			t.Errorf("Expected message to be 'OK', got %v", result.message)
		}
	})

	// Test case 2: Nil result
	t.Run("Nil result", func(t *testing.T) {
		_, err := toCommandResult(nil)

		if err == nil || !errors.Is(err, ErrMalformedResponseReceived) {
			t.Fatalf("Expected ErrMalformedResponseReceived, got %v", err)
		}
	})

	t.Run("Malformed response - wrong structure", func(t *testing.T) {
		input := []interface{}{"some-value", "invalid-code", "OK"}
		_, err := toCommandResult(input)

		if err == nil || !errors.Is(err, ErrMalformedResponseReceived) {
			t.Fatalf("Expected ErrMalformedResponseReceived, got %v", err)
		}
	})

	t.Run("Malformed response - wrong length", func(t *testing.T) {
		input := []interface{}{"some-value", int64(200)}
		_, err := toCommandResult(input)

		if err == nil || !errors.Is(err, ErrMalformedResponseReceived) {
			t.Fatalf("Expected ErrMalformedResponseReceived, got %v", err)
		}
	})

	t.Run("Malformed response - invalid types", func(t *testing.T) {
		input := []interface{}{123, int64(200), 456}
		_, err := toCommandResult(input)

		if err == nil || !errors.Is(err, ErrMalformedResponseReceived) {
			t.Fatalf("Expected ErrMalformedResponseReceived, got %v", err)
		}
	})
}
