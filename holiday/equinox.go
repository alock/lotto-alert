package holiday

import "time"

// SeasonalDates holds the equinox and solstice dates for a given year.
type SeasonalDates struct {
	SpringEquinox  time.Time
	SummerSolstice time.Time
	AutumnEquinox  time.Time
	WinterSolstice time.Time
}

// seasonalTable contains known equinox/solstice dates matching lottery data.
// Dates reflect what the lottery organization used, which may differ slightly
// from precise astronomical times due to timezone rounding.
var seasonalTable = map[int][4][2]int{
	// [spring month/day, summer month/day, autumn month/day, winter month/day]
	2023: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2024: {{3, 19}, {6, 21}, {9, 22}, {12, 21}},
	2025: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2026: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2027: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2028: {{3, 19}, {6, 20}, {9, 22}, {12, 21}},
	2029: {{3, 20}, {6, 20}, {9, 22}, {12, 21}},
	2030: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2031: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2032: {{3, 19}, {6, 20}, {9, 22}, {12, 21}},
	2033: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2034: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
	2035: {{3, 20}, {6, 21}, {9, 22}, {12, 21}},
}

// GetSeasonalDates returns the equinox/solstice dates for a given year.
// Falls back to typical dates if the year is not in the lookup table.
func GetSeasonalDates(year int) SeasonalDates {
	entry, ok := seasonalTable[year]
	if !ok {
		// Fallback to typical dates
		entry = [4][2]int{{3, 20}, {6, 21}, {9, 22}, {12, 21}}
	}
	return SeasonalDates{
		SpringEquinox:  time.Date(year, time.Month(entry[0][0]), entry[0][1], 0, 0, 0, 0, time.UTC),
		SummerSolstice: time.Date(year, time.Month(entry[1][0]), entry[1][1], 0, 0, 0, 0, time.UTC),
		AutumnEquinox:  time.Date(year, time.Month(entry[2][0]), entry[2][1], 0, 0, 0, 0, time.UTC),
		WinterSolstice: time.Date(year, time.Month(entry[3][0]), entry[3][1], 0, 0, 0, 0, time.UTC),
	}
}
