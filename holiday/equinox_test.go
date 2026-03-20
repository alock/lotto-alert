package holiday

import (
	"testing"
	"time"
)

func TestGetSeasonalDates(t *testing.T) {
	cases := []struct {
		year   int
		spring [2]int // month, day
		summer [2]int
		autumn [2]int
		winter [2]int
	}{
		{2023, [2]int{3, 20}, [2]int{6, 21}, [2]int{9, 22}, [2]int{12, 21}},
		{2024, [2]int{3, 19}, [2]int{6, 21}, [2]int{9, 22}, [2]int{12, 21}},
		{2025, [2]int{3, 20}, [2]int{6, 21}, [2]int{9, 22}, [2]int{12, 21}},
		{2026, [2]int{3, 20}, [2]int{6, 21}, [2]int{9, 22}, [2]int{12, 21}},
	}
	for _, tc := range cases {
		sd := GetSeasonalDates(tc.year)
		check := func(name string, got time.Time, month, day int) {
			want := time.Date(tc.year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			if !got.Equal(want) {
				t.Errorf("%d %s = %v, want %v", tc.year, name, got.Format("Jan 2"), want.Format("Jan 2"))
			}
		}
		check("Spring", sd.SpringEquinox, tc.spring[0], tc.spring[1])
		check("Summer", sd.SummerSolstice, tc.summer[0], tc.summer[1])
		check("Autumn", sd.AutumnEquinox, tc.autumn[0], tc.autumn[1])
		check("Winter", sd.WinterSolstice, tc.winter[0], tc.winter[1])
	}
}

func TestGetSeasonalDatesFallback(t *testing.T) {
	sd := GetSeasonalDates(2099)
	if sd.SpringEquinox.Month() != time.March || sd.SpringEquinox.Day() != 20 {
		t.Errorf("fallback Spring = %v, want Mar 20", sd.SpringEquinox.Format("Jan 2"))
	}
}
