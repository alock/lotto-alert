package prize

import (
	"testing"
	"time"
)

func TestForDateHoliday(t *testing.T) {
	// Christmas 2025
	info := ForDate(time.Date(2025, time.December, 25, 0, 0, 0, 0, time.UTC))
	if info.Amount != 250 {
		t.Errorf("Christmas amount = %d, want 250", info.Amount)
	}
	if info.Reason != "Christmas Day" {
		t.Errorf("Christmas reason = %q, want %q", info.Reason, "Christmas Day")
	}
}

func TestForDateSaturday(t *testing.T) {
	// Find a Saturday that isn't a holiday
	// Jan 3, 2026 is a Saturday
	info := ForDate(time.Date(2026, time.January, 3, 0, 0, 0, 0, time.UTC))
	if info.Amount != 50 {
		t.Errorf("Saturday amount = %d, want 50", info.Amount)
	}
	if info.Reason != "a Saturday" {
		t.Errorf("Saturday reason = %q, want %q", info.Reason, "a Saturday")
	}
}

func TestForDateWeekday(t *testing.T) {
	// A regular weekday with no holiday
	// Jan 5, 2026 is a Monday
	info := ForDate(time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC))
	if info.Amount != 30 {
		t.Errorf("weekday amount = %d, want 30", info.Amount)
	}
	if info.Reason != "" {
		t.Errorf("weekday reason = %q, want empty", info.Reason)
	}
}
