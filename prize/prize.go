package prize

import (
	"time"

	"github.com/alock/lotto-alert/holiday"
)

// Info holds prize amount and reason for a given date.
type Info struct {
	Amount int
	Reason string
}

// ForDate returns the prize info for a given date.
func ForDate(t time.Time) Info {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	holidays := holiday.ForYear(t.Year())
	if h, ok := holidays[date]; ok {
		return Info{Amount: h.Amount, Reason: h.Name}
	}
	if t.Weekday() == time.Saturday {
		return Info{Amount: 50, Reason: "a Saturday"}
	}
	return Info{Amount: 30, Reason: ""}
}
