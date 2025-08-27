package helpers

import (
	"time"
)

func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

func FormatTime(t time.Time) string {
	return t.Format("3:04 PM MST")
}
