package timeutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
)

func TestHasMonth(t *testing.T) {
	testCases := map[string]bool{
		time.January.String():   true,
		time.February.String():  true,
		time.March.String():     true,
		time.April.String():     true,
		time.May.String():       true,
		time.June.String():      true,
		time.July.String():      true,
		time.August.String():    true,
		time.September.String(): true,
		time.October.String():   true,
		time.November.String():  true,
		time.December.String():  true,
		"2st January 2009":      true,
		"mayday, mayday":        true,
		"next year":             false,
		"30th of feb":           false,
	}

	for in, want := range testCases {
		t.Run(in, func(t *testing.T) {
			if got := HasMonth(in); got != want {
				t.Errorf("HasMonth(%s) = %t, want %t", in, got, want)
			}
		})
	}
}

func TestParse(t *testing.T) {

}

func TestMonthNumber(t *testing.T) {
	testCases := map[string]int{
		"january":              1,
		time.December.String(): 12,
		"not a month":          0,
	}

	for in, want := range testCases {
		t.Run(in, func(t *testing.T) {
			if got := monthNumber(in); got != want {
				t.Errorf("monthNumber(%s) = %d, want %d", in, got, want)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := map[time.Time]string{
		time.Date(2021, time.March, 6, 0, 0, 0, 0, time.Local): "2021-03-06",
	}

	for in, want := range testCases {
		if got := String(in); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	}
}

func TestJSONTime_UnmarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		want time.Time
		err  error
	}{
		`"2021-03-06"`: {
			time.Date(2021, time.March, 6, 0, 0, 0, 0, time.Local),
			nil,
		},
	}

	for in, tc := range testCases {
		t.Run(in, func(t *testing.T) {

			var gotJT JSONTime
			if err := json.Unmarshal([]byte(in), &gotJT); err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("Unmarshal() err = %v, want %v", err, tc.err)
				}
				return
			}

			got := time.Time(gotJT)

			if diff := deep.Equal(got, tc.want); diff != nil {
				t.Errorf("Unmarshal() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestJSONTime_MarshalJSON(t *testing.T) {
	testCases := map[time.Time]string{
		time.Date(2021, time.March, 6, 0, 0, 0, 0, time.Local): `"2021-03-06"`,
	}

	for in, want := range testCases {
		t.Run(in.Format(Format), func(t *testing.T) {
			gotBytes, err := json.Marshal(JSONTime(in))
			if err != nil {
				t.Fatalf("Marshal() err = %v, want %v", err, nil)
			}

			if got := string(gotBytes); got != want {
				t.Errorf("Marshal() = %v, want %v", got, want)
			}
		})
	}
}

func TestDaysBetween(t *testing.T) {
	testCases := []struct {
		start string
		end   string
		want  []string
		err   error
	}{
		{
			"2021-03-01", "2021-03-04",
			[]string{"2021-03-02", "2021-03-03"},
			nil,
		}, {
			"2021-03-01", "2021-03-02",
			nil,
			nil,
		}, {
			"2021-03-02", "2021-03-01",
			nil,
			ErrInvalidRange,
		},
	}

	for _, tc := range testCases {
		start, err := time.Parse(Format, tc.start)
		if err != nil {
			t.Fatalf("can't parse start = %v", err)
		}

		end, err := time.Parse(Format, tc.end)
		if err != nil {
			t.Fatalf("can't parse end = %v", err)
		}

		want := make([]time.Time, 0, len(tc.want))
		for _, w := range tc.want {
			wt, err := time.Parse(Format, w)
			if err != nil {
				t.Fatalf("can't parse test want = %v", err)
			}
			want = append(want, wt)
		}

		t.Run(fmt.Sprintf("%s~%s", tc.start, tc.end), func(t *testing.T) {
			got, err := DaysBetween(start, end)
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("DaysBetween(%s, %s) err = %v, want %v",
						start, end, err, tc.err)
				}
				return
			}
			if diff := deep.Equal(got, want); diff != nil {
				t.Errorf("DaysBetween(%s, %s) = %v, want %v, diff = %v",
					start, end, got, tc.want, diff)
			}
		})
	}
}
