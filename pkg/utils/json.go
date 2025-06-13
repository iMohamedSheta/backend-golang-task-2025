package utils

import (
	"reflect"
)

// GetJSONKeys returns the JSON keys of the given struct
func GetJSONKeys(i any) map[string]bool {
	keys := make(map[string]bool)

	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			continue
		}

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		name := tag
		if idx := indexComma(tag); idx != -1 {
			name = tag[:idx]
		}

		keys[name] = true
	}

	return keys
}

func indexComma(tag string) int {
	for i, r := range tag {
		if r == ',' {
			return i
		}
	}
	return -1
}
