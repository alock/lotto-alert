package holiday

import (
	"testing"
	"time"
)

func TestEaster(t *testing.T) {
	cases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2023, time.April, 9},
		{2024, time.March, 31},
		{2025, time.April, 20},
		{2026, time.April, 5},
		{2027, time.March, 28},
		{2028, time.April, 16},
		{2029, time.April, 1},
		{2030, time.April, 21},
		{2031, time.April, 13},
		{2032, time.March, 28},
		{2033, time.April, 17},
		{2034, time.April, 9},
		{2035, time.March, 25},
	}
	for _, tc := range cases {
		got := Easter(tc.year)
		want := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("Easter(%d) = %v, want %v", tc.year, got.Format("Jan 2"), want.Format("Jan 2"))
		}
	}
}
