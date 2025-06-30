package utils

import (
	"fmt"
	"taskgo/pkg/enums"
)

func Print(v ...any) {
	for _, val := range v {
		fmt.Println(val)
	}
}

func PrintErr(v ...any) {
	for _, val := range v {
		fmt.Println(enums.Red.Value() + fmt.Sprint(val) + enums.Reset.Value())
	}
}

func PrintSuccess(v ...any) {
	for _, val := range v {
		fmt.Println(enums.Green.Value() + fmt.Sprint(val) + enums.Reset.Value())
	}
}

func PrintWarning(v ...any) {
	for _, val := range v {
		fmt.Println(enums.Yellow.Value() + fmt.Sprint(val) + enums.Reset.Value())
	}
}

func PrintInfo(v ...any) {
	for _, val := range v {
		fmt.Println(enums.Blue.Value() + fmt.Sprint(val) + enums.Reset.Value())
	}
}
