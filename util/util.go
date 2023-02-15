package util

import (
	"time"
)

// GetStringOfDate simple helper to always format the date to simple M/D/YYYY
func GetStringOfDate(t time.Time) string {
	return t.Format("01/02/06")
}

func TruncateToDayValue(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
