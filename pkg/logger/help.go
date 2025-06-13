package logger

import (
	"fmt"
	"time"
)

func LogTimer(title string, start time.Time) {
	Log().Info(fmt.Sprintf("%s took %s", title, time.Since(start)))
}
