package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"taskgo/pkg/enums"
)

func DD(v ...any) {
	Dump(v...)
	Die()
}

func Dump(v ...any) {
	for _, val := range v {
		fmt.Println(formatOutput(val))
	}
}

func Die() {
	os.Exit(1)
}

func formatOutput(value any) string {
	red := enums.Green.Value()
	reset := enums.Reset.Value()
	blackBG := enums.BG_Black.Value()

	// Handle different types
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%s%s%s%s", blackBG, red, v, reset)
	case error:
		return fmt.Sprintf("%s%s%s%s", blackBG, red, v.Error(), reset)
	default:
		// Use JSON for complex structures
		jsonData, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%s%s%#v%s", blackBG, red, value, reset)
		}
		return fmt.Sprintf("%s%s%s%s", blackBG, red, jsonData, reset)
	}
}
