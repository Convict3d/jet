package jet

import (
	"fmt"
	"reflect"
	"strings"
)

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v)
}

func MarshalFlattened(v interface{}) ([]byte, error) {
	return marshalFlattened(v)
}

func MarshalNormalized(v interface{}) ([]byte, error) {
	return marshalNormalized(v)
}

func marshalNormalized(v interface{}) ([]byte, error) {
	genericData, err := encode(v)
	if err != nil {
		return nil, err
	}

	formattedBytes, err := formatNormalized(genericData)
	if err != nil {
		return nil, err
	}

	return formattedBytes, nil
}

func marshalFlattened(v interface{}) ([]byte, error) {
	genericData, err := encode(v)
	if err != nil {
		return nil, err
	}

	formattedBytes, err := formatFlattened(genericData)
	if err != nil {
		return nil, err
	}

	return formattedBytes, nil
}

func Unmarshal(data []byte, v interface{}) error {
	return unmarshal(data, v)
}

func marshal(v interface{}) ([]byte, error) {

	genericData, err := encode(v)
	if err != nil {
		return nil, err
	}

	formattedBytes, err := format(genericData)
	if err != nil {
		return nil, err
	}

	return formattedBytes, nil
}

func encode(v interface{}) (interface{}, error) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		resultMap := make(map[string]interface{})
		t := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := t.Field(i)
			fieldValue := val.Field(i)

			tagName := field.Tag.Get("jet")
			if tagName == "" {
				tagName = strings.ToLower(field.Name)
			}
			if tagName == "-" {
				continue // Skip this field
			}

			encodedValue, err := encode(fieldValue.Interface())
			if err != nil {
				return nil, err
			}
			resultMap[tagName] = encodedValue
		}
		return resultMap, nil
	case reflect.Slice:
		resultSlice := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			encodedValue, err := encode(val.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			resultSlice[i] = encodedValue
		}
		return resultSlice, nil
	case reflect.Map:
		resultMap := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			mapValue := val.MapIndex(key)
			encodedValue, err := encode(mapValue.Interface())
			if err != nil {
				return nil, err
			}
			resultMap[key.String()] = encodedValue
		}
		return resultMap, nil
	case reflect.String, reflect.Int, reflect.Int64, reflect.Int32, reflect.Float64, reflect.Float32, reflect.Bool:
		return val.Interface(), nil
	default:
		return nil, fmt.Errorf("jet: unsupported type for marshaling %s", val.Kind())
	}
}

func unmarshal(data []byte, v interface{}) error {
	// Implementation of unmarshaling logic goes here
	return nil
}
