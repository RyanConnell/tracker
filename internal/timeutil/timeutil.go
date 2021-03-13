// Package timeutil helps with time parsing and human readable times.
package timeutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	// precompute all the
	months = make(map[string]struct{})
	for m := 1; m <= 12; m++ {
		months[time.Month(m).String()] = struct{}{}
	}
}

const Format = "2006-01-02"

const Day = 24 * time.Hour

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrInvalidTime  = Error("timeutil: unable to parse time")
	ErrInvalidRange = Error("timeutil: invalid time range")
)

// months contains precomputed map of months so they are quick to access.
var months map[string]struct{}

// HasMonth checks is there a month in the string, which indicates that the
// string is possibly a date. The month must be a whole word.
func HasMonth(str string) bool {
	for _, word := range strings.Fields(str) {
		word := strings.Map(func(r rune) rune {
			if strings.ContainsRune(".,:;", r) {
				return -1
			} else {
				return r
			}
		}, strings.Title(strings.ToLower(word)))

		if _, exists := months[word]; exists {
			return true
		}
	}

	return false
}

var timeRegexp = regexp.MustCompile(`([a-zA-Z]+)[^0-9]+([0-9]+)[^0-9]+([0-9]+)`)

// Parse the time
// TODO: If this implements a valid time RFC, we should convert to that instead.
//       Maybe using time.Parse(Format, str) would suffice.
func Parse(str string) (t time.Time, err error) {
	var day, month, year int

	if matches := timeRegexp.FindStringSubmatch(str); len(matches) >= 3 {
		if day, err = strconv.Atoi(matches[2]); err != nil {
			return t, fmt.Errorf("unable to parse day %s: %w",
				matches[2], err)
		}
		if month = monthNumber(matches[1]); month == 0 {
			return t, fmt.Errorf("invalid month: %s", matches[1])
		}
		if year, err = strconv.Atoi(matches[3]); err != nil {
			return t, fmt.Errorf("unable to parse year %s: %w",
				matches[3], err)
		}
	} else {
		return t, ErrInvalidTime
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// monthNumber returns the number of the month
func monthNumber(m string) int {
	for i := 1; i <= 12; i++ {
		if strings.EqualFold(time.Month(i).String(), m) {
			return i
		}
	}

	return 0
}

// String converts the time to a string in a consistent format
func String(t time.Time) string {
	return t.Format(Format)
}

// JSONTime allows to convert time to and from JSON
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format(Format))), nil
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.TrimPrefix(strings.TrimSuffix(string(b), "\""), "\"")

	n, err := time.Parse(Format, s)
	if err != nil {
		return err
	}

	*t = JSONTime(n)
	return nil
}

// DaysBetween returns all the posible days between the two time ranges.
func DaysBetween(start, end time.Time) ([]time.Time, error) {
	if end.Before(start) {
		return nil, ErrInvalidRange
	}

	days := make([]time.Time, 0)

	next := start.Add(Day)
	for next.Before(end) {
		days = append(days, next)
		next = next.Add(Day)
	}

	return days, nil
}
