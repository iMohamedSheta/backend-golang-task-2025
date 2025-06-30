package reflect

import (
	"reflect"
	"strings"
)

// GetJsonFieldName returns the JSON field name for a struct field
// Returns empty string if field not found or should be skipped (json:"-")
func GetJsonFieldName(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return ""
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName // Use original field name if no json tag
	}

	// Handle json:"field_name,omitempty" format
	if idx := strings.Index(jsonTag, ","); idx != -1 {
		jsonTag = jsonTag[:idx]
	}

	// If json tag is "-", return empty string to indicate skip
	if jsonTag == "-" {
		return ""
	}

	return jsonTag
}
