package scrape

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestParsePaLottoResults(t *testing.T) {
	data, err := os.ReadFile("testdata/sample.html")
	if err != nil {
		t.Fatalf("failed to read sample HTML: %v", err)
	}
	results := parsePaLottoResults(string(data))

	if len(results) == 0 {
		t.Fatal("parsed zero results from sample HTML")
	}
	if len(results)%2 != 0 {
		t.Fatalf("expected even number of results (date/number pairs), got %d", len(results))
	}

	pairs := len(results) / 2
	t.Logf("parsed %d date/number pairs", pairs)

	// Validate every pair is a valid date + 0-999 number
	for i := 0; i < len(results); i += 2 {
		dateStr := results[i]
		numStr := results[i+1]

		_, err := time.ParseInLocation("1/2/2006", dateStr, time.Local)
		if err != nil {
			t.Errorf("pair %d: invalid date %q: %v", i/2, dateStr, err)
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			t.Errorf("pair %d: invalid number %q: %v", i/2, numStr, err)
		}
		if num < 0 || num > 999 {
			t.Errorf("pair %d: number %d out of Pick 3 range (0-999)", i/2, num)
		}
	}

	// Page returns reverse chronological order; last entry should be 1/1/2026
	lastDate := results[len(results)-2]
	d, _ := time.ParseInLocation("1/2/2006", lastDate, time.Local)
	if d.Year() != 2026 || d.Month() != time.January || d.Day() != 1 {
		t.Errorf("last entry date = %s, expected 1/1/2026", lastDate)
	}
	lastNum := results[len(results)-1]
	if lastNum != "328" {
		t.Errorf("last entry number = %s, expected 328 (known Jan 1 2026 result)", lastNum)
	}
}
