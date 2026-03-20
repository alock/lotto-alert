package holiday

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

type fixtureEntry struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}

func TestHolidaysMatchHistoricalData(t *testing.T) {
	for _, year := range []int{2023, 2024, 2025, 2026} {
		t.Run(fmt.Sprintf("%d", year), func(t *testing.T) {
			data, err := os.ReadFile(fmt.Sprintf("testdata/%d.json", year))
			if err != nil {
				t.Fatalf("failed to read fixture: %v", err)
			}
			var expected map[string]fixtureEntry
			if err := json.Unmarshal(data, &expected); err != nil {
				t.Fatalf("failed to parse fixture: %v", err)
			}

			computed := ForYear(year)

			// Build map keyed by date string for comparison
			computedMap := make(map[string]fixtureEntry)
			for date, h := range computed {
				key := date.Format("01/02/06")
				computedMap[key] = fixtureEntry{Amount: h.Amount, Reason: h.Name}
			}

			// Check every expected entry exists in computed
			for key, exp := range expected {
				comp, ok := computedMap[key]
				if !ok {
					t.Errorf("missing holiday %s (%s)", key, exp.Reason)
					continue
				}
				if comp.Amount != exp.Amount {
					t.Errorf("%s (%s): amount = %d, want %d", key, exp.Reason, comp.Amount, exp.Amount)
				}
				if comp.Reason != exp.Reason {
					t.Errorf("%s: reason = %q, want %q", key, comp.Reason, exp.Reason)
				}
			}

			// Check for extra computed holidays not in fixture
			for key, comp := range computedMap {
				if _, ok := expected[key]; !ok {
					t.Errorf("extra holiday %s (%s, $%d)", key, comp.Reason, comp.Amount)
				}
			}
		})
	}
}

func TestForYearCollisionResolution(t *testing.T) {
	// 2026: Father's Day (Jun 21, $125) collides with Summer Solstice (Jun 21, $50)
	holidays := ForYear(2026)
	jun21 := time.Date(2026, time.June, 21, 0, 0, 0, 0, time.UTC)
	h, ok := holidays[jun21]
	if !ok {
		t.Fatal("expected a holiday on Jun 21, 2026")
	}
	if h.Name != "Father's Day" {
		t.Errorf("Jun 21 2026: got %q, want Father's Day (collision resolution)", h.Name)
	}
	if h.Amount != 125 {
		t.Errorf("Jun 21 2026: amount = %d, want 125", h.Amount)
	}
}
