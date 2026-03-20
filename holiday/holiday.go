package holiday

import "time"

// Holiday represents a special date with a prize amount.
type Holiday struct {
	Name      string
	Amount    int
	DateFunc  func(year int) time.Time
	StartYear int // 0 = always active
	EndYear   int // 0 = always active
}

// Date helper constructors

func fixedDate(month time.Month, day int) func(int) time.Time {
	return func(year int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
}

func nthWeekday(month time.Month, weekday time.Weekday, n int) func(int) time.Time {
	return func(year int) time.Time {
		first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		offset := (int(weekday) - int(first.Weekday()) + 7) % 7
		day := 1 + offset + (n-1)*7
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}
}

func lastWeekday(month time.Month, weekday time.Weekday) func(int) time.Time {
	return func(year int) time.Time {
		// Start from the last day of the month and work backwards
		last := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
		offset := (int(last.Weekday()) - int(weekday) + 7) % 7
		return last.AddDate(0, 0, -offset)
	}
}

func easterOffset(days int) func(int) time.Time {
	return func(year int) time.Time {
		return Easter(year).AddDate(0, 0, days)
	}
}

// electionDay returns the first Tuesday after the first Monday in November.
func electionDay(year int) time.Time {
	firstMonday := nthWeekday(time.November, time.Monday, 1)(year)
	return firstMonday.AddDate(0, 0, 1)
}

// registry contains all known holidays.
var registry = []Holiday{
	{Name: "New Year's Day", Amount: 200, DateFunc: fixedDate(time.January, 1)},
	{Name: "Martin Luther King Jr Day", Amount: 50, DateFunc: nthWeekday(time.January, time.Monday, 3)},
	{Name: "Groundhog Day", Amount: 50, DateFunc: fixedDate(time.February, 2)},
	{Name: "Lincoln's Birthday", Amount: 50, DateFunc: fixedDate(time.February, 12)},
	{Name: "Valentine's Day", Amount: 100, DateFunc: fixedDate(time.February, 14)},
	{Name: "President's Day", Amount: 50, DateFunc: nthWeekday(time.February, time.Monday, 3)},
	{Name: "Washington's Birthday", Amount: 40, DateFunc: fixedDate(time.February, 22), EndYear: 2023},
	{Name: "St Patrick's Day", Amount: 50, DateFunc: fixedDate(time.March, 17)},
	{Name: "the First Day of Spring", Amount: 50, DateFunc: func(year int) time.Time {
		return GetSeasonalDates(year).SpringEquinox
	}},
	{Name: "April Fool's Day (for real)", Amount: 60, DateFunc: fixedDate(time.April, 1)},
	{Name: "Palm Sunday", Amount: 80, DateFunc: easterOffset(-7)},
	{Name: "Good Friday", Amount: 75, DateFunc: easterOffset(-2)},
	{Name: "Easter", Amount: 100, DateFunc: easterOffset(0)},
	{Name: "John Muir's Birthday", Amount: 50, DateFunc: fixedDate(time.April, 21)},
	{Name: "Earth Day", Amount: 100, DateFunc: fixedDate(time.April, 22)},
	{Name: "John Audubon's Birthday", Amount: 50, DateFunc: fixedDate(time.April, 26)},
	{Name: "Mother's Day", Amount: 125, DateFunc: nthWeekday(time.May, time.Sunday, 2)},
	{Name: "Memorial Day", Amount: 100, DateFunc: lastWeekday(time.May, time.Monday)},
	{Name: "Flag Day", Amount: 50, DateFunc: fixedDate(time.June, 14), EndYear: 2023},
	{Name: "Juneteenth", Amount: 50, DateFunc: fixedDate(time.June, 19), StartYear: 2026},
	{Name: "Father's Day", Amount: 125, DateFunc: nthWeekday(time.June, time.Sunday, 3)},
	{Name: "the First Day of Summer", Amount: 50, DateFunc: func(year int) time.Time {
		return GetSeasonalDates(year).SummerSolstice
	}},
	{Name: "Independence Day", Amount: 150, DateFunc: fixedDate(time.July, 4)},
	{Name: "Labor Day", Amount: 100, DateFunc: nthWeekday(time.September, time.Monday, 1)},
	{Name: "Patriot's Day", Amount: 50, DateFunc: fixedDate(time.September, 11)},
	{Name: "the First Day of Autumn", Amount: 65, DateFunc: func(year int) time.Time {
		return GetSeasonalDates(year).AutumnEquinox
	}},
	{Name: "Indigenous Peoples Day", Amount: 50, DateFunc: nthWeekday(time.October, time.Monday, 2)},
	{Name: "Halloween", Amount: 75, DateFunc: fixedDate(time.October, 31)},
	{Name: "Election Day", Amount: 50, DateFunc: electionDay},
	{Name: "Veteran's Day", Amount: 75, DateFunc: fixedDate(time.November, 11)},
	{Name: "Thanksgiving", Amount: 200, DateFunc: nthWeekday(time.November, time.Thursday, 4)},
	{Name: "Pearl Harbor Day", Amount: 50, DateFunc: fixedDate(time.December, 7)},
	{Name: "the day Wildlife Works Inc started", Amount: 75, DateFunc: fixedDate(time.December, 9)},
	{Name: "the First Day of Winter", Amount: 75, DateFunc: func(year int) time.Time {
		return GetSeasonalDates(year).WinterSolstice
	}},
	{Name: "Christmas Eve", Amount: 150, DateFunc: fixedDate(time.December, 24)},
	{Name: "Christmas Day", Amount: 250, DateFunc: fixedDate(time.December, 25)},
	{Name: "New Year's Eve", Amount: 150, DateFunc: fixedDate(time.December, 31)},
}

// ForYear computes all holidays for a given year.
// When two holidays fall on the same date, the one with the highest prize wins.
func ForYear(year int) map[time.Time]Holiday {
	result := make(map[time.Time]Holiday)
	for _, h := range registry {
		if h.StartYear > 0 && year < h.StartYear {
			continue
		}
		if h.EndYear > 0 && year > h.EndYear {
			continue
		}
		date := h.DateFunc(year)
		if existing, ok := result[date]; ok {
			if h.Amount > existing.Amount {
				result[date] = h
			}
		} else {
			result[date] = h
		}
	}
	return result
}
