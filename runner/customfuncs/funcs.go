package customfuncs

import "math"

// IsTruthy mirrors JavaScript's !! coercion rules.
// Arrays and objects are always truthy; zero, empty string, nil, and false are not.
func IsTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case string:
		return val != ""
	default:
		return true
	}
}

// EncodeBitmask converts a slice of values to a bitmask integer.
// Each truthy value at index i contributes 2^i to the result.
func EncodeBitmask(values interface{}) (float64, error) {
	arr, ok := values.([]interface{})
	if !ok {
		return 0, nil
	}
	result := 0.0
	for i, v := range arr {
		if IsTruthy(v) {
			result += math.Pow(2, float64(i))
		}
	}
	return result, nil
}

// DecodeBitmask converts a bitmask integer back to a map using the provided keys.
// Null keys are skipped; each present key maps to whether its bit is set.
func DecodeBitmask(data float64, keys interface{}) (interface{}, error) {
	arr, ok := keys.([]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}
	result := make(map[string]interface{}, len(arr))
	intData := int64(data)
	for i, key := range arr {
		if key == nil {
			continue
		}
		keyStr, ok := key.(string)
		if !ok {
			continue
		}
		result[keyStr] = (intData & (1 << uint(i))) != 0
	}
	return result, nil
}

// Bitmask returns whether the bit at the given index is set in data.
func Bitmask(data float64, index float64) bool {
	return (int64(data) & (1 << uint(index))) != 0
}
