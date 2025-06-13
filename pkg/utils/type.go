package utils

import (
	"strings"
)

// Convert value type (array of strings) to CSV string
func ToCSV(val any, defaultValue string) string {
	switch v := val.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ", ")
	default:
		return defaultValue
	}
}

func ToArrayOfStrings(val any, defaultValue []string) []string {
	switch v := val.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return defaultValue
	}
}
