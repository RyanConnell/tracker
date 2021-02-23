package date

import (
	"testing"
)

func TestToDate(t *testing.T) {
	var testCases = []struct {
		date         string
		expectedDate *Date
	}{
		{"January 1, 2010", &Date{1, 1, 2010}},
		{"January 5, 1995", &Date{5, 1, 1995}},
		{"August 24, 2040 (2010-08-24) (test)", &Date{24, 8, 2040}},
	}

	for _, tc := range testCases {
		t.Run(tc.date, func(t *testing.T) {
			date, err := ToDate(tc.date)
			if err != nil {
				t.Fatalf("Error occured creating date: %v\n", err)
			}
			if date.CompareTo(tc.expectedDate) != 0 {
				t.Fatalf("Dates were not equal: %#v != %#v\n", date, tc.expectedDate)
			}
		})
	}
}
