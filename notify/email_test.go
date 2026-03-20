package notify

import (
	"testing"
	"time"

	"github.com/alock/lotto-alert/prize"
)

func TestGetMessage(t *testing.T) {
	cases := []struct {
		name   string
		date   time.Time
		number int
		prize  prize.Info
		want   string
	}{
		{
			name:   "holiday",
			date:   time.Date(2026, time.December, 25, 0, 0, 0, 0, time.UTC),
			number: 123,
			prize:  prize.Info{Amount: 250, Reason: "Christmas Day"},
			want:   "Congrats! On December 25 the PICK 3 Evening Number was 123 and you won $250 because it is Christmas Day.",
		},
		{
			name:   "no reason",
			date:   time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC),
			number: 7,
			prize:  prize.Info{Amount: 30, Reason: ""},
			want:   "Congrats! On January 5 the PICK 3 Evening Number was 007 and you won $30.",
		},
		{
			name:   "saturday",
			date:   time.Date(2026, time.January, 3, 0, 0, 0, 0, time.UTC),
			number: 42,
			prize:  prize.Info{Amount: 50, Reason: "a Saturday"},
			want:   "Congrats! On January 3 the PICK 3 Evening Number was 042 and you won $50 because it is a Saturday.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetMessage(tc.date, tc.number, tc.prize)
			if got != tc.want {
				t.Errorf("got:  %s\nwant: %s", got, tc.want)
			}
		})
	}
}
