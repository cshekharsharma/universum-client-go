package universum

import (
	"bufio"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Encode converts a Go data type into its corresponding RESP3 encoded string format.
// It supports encoding basic types (string, integers, floats), composite types (slices, maps),
// booleans, nil values, custom error messages, and even specific custom types like
// *storage.ScalarRecord. The function uses type assertion and reflection to determine
// the input type and formats it accordingly.
//
// Parameters:
//   - value interface{}: The value to be encoded into RESP3 format. This could be any supported
//     Go data type, including custom types that have a defined encoding pattern.
//
// Returns:
//   - string: The RESP3 encoded string representation of the input `value`.
//   - error: An error is returned if the value type is not supported for encoding or
//     if any issue arises during the encoding process.
//
// Supported Types:
//   - Basic types: Encodes strings, integers, and floats with respective RESP3 prefixes.
//   - Composite types: Encodes slices as RESP3 arrays and maps as RESP3 maps, recursively encoding
//     their elements.
//   - Booleans: Encodes true and false as RESP3 booleans.
//   - Nil: Encodes a nil value as RESP3 Null.
//   - Errors: Encodes Go error types as RESP3 error messages.
//   - *storage.ScalarRecord: Encodes custom struct types by converting them into a generic map
//     and then encoding this map.
func encodeResp(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return "$" + strconv.Itoa(len(v)) + "\r\n" + v + "\r\n", nil

	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return ":" + fmt.Sprintf("%d", v) + "\r\n", nil

	case float32, float64:
		return "," + fmt.Sprintf("%f", v) + "\r\n", nil

	case []interface{}:
		resp := "*" + strconv.Itoa(len(v)) + "\r\n"
		for _, elem := range v {
			encodedElem, err := encodeResp(elem)
			if err != nil {
				return "", err
			}

			resp += encodedElem
		}
		return resp, nil

	case []string:
		resp := "*" + strconv.Itoa(len(v)) + "\r\n"
		for _, elem := range v {
			encodedElem, err := encodeResp(elem)
			if err != nil {
				return "", err
			}

			resp += encodedElem
		}
		return resp, nil

	case bool:
		if v {
			return "#t\r\n", nil
		}
		return "#f\r\n", nil

	case nil:
		return "_\r\n", nil

	case error:
		return "-" + v.Error() + "\r\n", nil

	case map[string]interface{}:
		resp := "%" + strconv.Itoa(len(v)*2) + "\r\n"
		for kx, vx := range v {
			resp += "+" + kx + "\r\n"
			valueStr, err := encodeResp(vx)
			if err != nil {
				return "", err
			}
			resp += valueStr
		}
		return resp, nil

	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

// Decode reads from the provided bufio.Reader and interprets the next RESP3 data type,
// returning the parsed value as an interface{}. RESP3 protocol supports various data types
// like Simple Strings, Errors, Integers, Floats, Bulk Strings, Arrays, Booleans, Maps, and Nulls.
// This function is capable of decoding these types based on the initial byte that indicates
// the data type, followed by the data itself.
//
// Parameters:
//   - reader *bufio.Reader: A pointer to a bufio.Reader from which the data will be read. It is
//     expected that the reader is already initialized and points to a source of RESP3 formatted data.
//
// Returns:
//
//   - interface{}: The decoded data from the reader. The actual type of the returned value can be
//     one of several Go types depending on the RESP3 data type encountered. This could be a string
//     for Simple Strings and Bulk Strings, error for RESP3 Errors, int64 for Integers, float64 for
//     Floats, []interface{} for Arrays, map[string]interface{} for Maps, bool for Booleans, or nil
//     for Nulls.
//
//     Sample input string: "*5\r\n$4\r\nMSET\r\n$4\r\nkey1\r\n$16\r\nvalue1 dash dash\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n"
func decodeResp(reader *bufio.Reader) (interface{}, error) {
	dataType, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	switch dataType {
	case '+': // Simple String
		line, _, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		return string(line), nil

	case '-': // Error
		line, _, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		return errors.New(string(line)), nil

	case ':': // Integer
		line, _, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		xint, castErr := strconv.Atoi(string(line))
		if castErr != nil {
			return xint, castErr
		}
		return int64(xint), nil

	case ',': // Float
		line, _, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		xfloat, castErr := strconv.ParseFloat(string(line), 64)
		if castErr != nil {
			return xfloat, castErr
		}
		return float64(xfloat), nil

	case '$': // Bulk String
		lengthStr, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		length, err := strconv.Atoi(string(lengthStr))
		if err != nil {
			return nil, err
		}

		if length == -1 {
			return nil, nil // Null bulk string
		}

		value := make([]byte, length)
		_, err = reader.Read(value)
		if err != nil {
			return nil, err
		}
		_, err = reader.Discard(2)

		if err != nil {
			return nil, err
		}
		return string(value), nil

	case '*': // Array
		countStr, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		count, err := strconv.Atoi(string(countStr))
		if err != nil {
			return nil, err
		}

		if count == -1 {
			return nil, nil // Null array
		}

		array := make([]interface{}, count)

		for i := 0; i < count; i++ {
			element, err := decodeResp(reader)

			if err != nil {
				return nil, err
			}

			array[i] = element
		}
		return array, nil

	case '#':
		b, err := reader.ReadByte()
		if err != nil {
			return false, err
		}

		reader.Discard(2)
		return b == 't', nil

	case '%': // Map of interface{}
		line, _, _ := reader.ReadLine()
		size, _ := strconv.Atoi(string(line))
		resultMap := make(map[string]interface{}, size/2)
		for i := 0; i < size; i += 2 {
			key, _ := decodeResp(reader)
			value, _ := decodeResp(reader)
			resultMap[key.(string)] = value
		}
		return resultMap, nil

	case '_':
		reader.Discard(2) // Discard the trailing \r\n
		return nil, nil

	default:
		return nil, fmt.Errorf("unsupported data type: %v", dataType)
	}
}
