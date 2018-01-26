package common

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Date struct {
	Day   int
	Month int
	Year  int
}

type NullDate struct {
	Date  Date
	Valid bool
}

var months = []string{
	"january", "february", "march", "april", "may", "june", "july",
	"august", "september", "october", "november", "december",
}

var dateSuffixes = []string{
	"th", "st", "nd", "rd", "th", "th", "th", "th", "th", "th",
}

var daysPerMonth = []int{
	31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31,
}

var weekdays = []string{
	"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
}

func (d *Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d *Date) FancyString() string {
	return fmt.Sprintf("%s %d%s, %4d", strings.Title(months[d.Month-1]),
		d.Day, d.DaySuffix(), d.Year)
}

func (d *Date) DaySuffix() string {
	if d.Day >= 10 && d.Day <= 20 {
		return "th"
	}
	return dateSuffixes[d.Day%10]
}

func (d *Date) CalendarString() string {
	return fmt.Sprintf("%s. %s. %d%s", strings.Title(d.WeekdayString())[:3],
		strings.Title(months[d.Month-1])[:3], d.Day, d.DaySuffix())
}

func (d *Date) isEmpty() bool {
	return d.Day == 0 && d.Month == 0 && d.Year == 0
}

func (d *Date) ToTime() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

func (d *Date) CompareTo(date *Date) int {
	if date == nil {
		return 1
	}
	if d == nil {
		return -1
	}

	if d.Year == date.Year {
		if d.Month == date.Month {
			if d.Day == date.Day {
				return 0
			}
			if d.Day > date.Day {
				return 1
			}
		}
		if d.Month > date.Month {
			return 1
		}
	}
	if d.Year > date.Year {
		return 1
	}
	return -1
}

func (d *Date) Weekday() int {
	return int(d.ToTime().Weekday())
}

func (d *Date) WeekdayString() string {
	return weekdays[d.Weekday()]
}

// Returns the Date that takes place {i} days before {d}.
func (d *Date) Minus(i int) *Date {
	date := &Date{d.Day, d.Month, d.Year}
	date.Day -= i

	for date.Day <= 0 {
		date.Month -= 1
		if date.Month == 0 {
			date.Month = 12
			date.Year -= 1
		}
		date.Day += DaysPerMonth(date.Month, date.Year)
	}

	return date
}

// Returns the Date that takes place {i} days after {d}.
func (d *Date) Plus(i int) *Date {
	date := &Date{d.Day, d.Month, d.Year}
	date.Day += i

	for date.Day > DaysPerMonth(date.Month, date.Year) {
		date.Day -= DaysPerMonth(date.Month, date.Year)
		date.Month += 1
		if date.Month == 13 {
			date.Month = 1
			date.Year++
		}
	}

	return date
}

func DateFromStr(str string) (Date, error) {
	date := &NullDate{}
	if err := date.FromStr(str); err != nil {
		return Date{}, err
	}
	if !date.Valid {
		return Date{}, fmt.Errorf("Invalid Date: %s", str)
	}
	return date.Date, nil
}

func (nd *NullDate) Scan(value interface{}) error {
	if val, ok := value.([]uint8); ok {
		err := nd.FromStr(string(val))
		if err != nil {
			return err
		}
		nd.Valid = true
	} else {
		return fmt.Errorf("Unsupported type given for NullDate: %T", value)
	}
	return nil
}

func (nd NullDate) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return fmt.Sprintf("%04d-%02d-%02d", nd.Date.Year, nd.Date.Month, nd.Date.Day), nil
}

func (nd *NullDate) FromStr(str string) error {
	values := strings.Split(str, "-")
	var err error
	if nd.Date.Year, err = StringToInt(values[0]); err != nil {
		return err
	}
	if nd.Date.Month, err = StringToInt(values[1]); err != nil {
		return err
	}
	if nd.Date.Day, err = StringToInt(values[2]); err != nil {
		return err
	}
	nd.Valid = nd.isValid()
	return nil
}

func (nd *NullDate) isValid() bool {
	// Add a rule for years here too?
	if nd.Date.Month <= 0 || nd.Date.Month > 12 {
		return false
	}
	if nd.Date.Day <= 0 || nd.Date.Day > DaysPerMonth(nd.Date.Month, nd.Date.Year) {
		return false
	}
	return true
}

func IsDate(str string) bool {
	for _, month := range months {
		if strings.Contains(strings.ToLower(str), month) {
			return true
		}
	}
	return false
}

func ToDate(str string) (*Date, error) {
	date := &Date{}
	regex, err := regexp.Compile(`([a-zA-Z]+)[^0-9]+([0-9]+)[^0-9]+([0-9]+)`)
	matches := regex.FindStringSubmatch(str)
	if len(matches) >= 3 {
		if date.Day, err = StringToInt(matches[2]); err != nil {
			return nil, err
		}
		if date.Month = matchMonth(matches[1]); date.Month == 0 {
			return nil, fmt.Errorf("Invalid month: %s", matches[1])
		}
		if date.Year, err = StringToInt(matches[3]); err != nil {
			return nil, err
		}
	} else {
		return date, fmt.Errorf("Unable to match regexp against str")
	}
	return date, nil
}

func matchMonth(str string) int {
	for i, month := range months {
		if strings.ToLower(str) == month {
			return i + 1
		}
	}
	return 0
}

func DatesInRange(start, end Date) ([]*Date, error) {
	if start.CompareTo(&end) == 1 {
		return nil, fmt.Errorf("Starting Date takes place before End Date.")
	}

	years := end.Year - start.Year + 1
	months := (12 * years) - start.Month - (12 - end.Month) + 1

	dateRange := make([]*Date, 0)

	// Need to handle case where first month = final month
	if start.Year == end.Year && start.Month == end.Month {
		for i := start.Day; i <= end.Day; i++ {
			dateRange = append(dateRange, &Date{i, start.Month, start.Year})
		}
		return dateRange, nil
	}

	// Add first month
	for i := start.Day; i <= DaysPerMonth(start.Month, start.Year); i++ {
		dateRange = append(dateRange, &Date{i, start.Month, start.Year})
	}

	// Add intermediate months
	for i := 1; i < months-1; i++ {
		month := start.Month + (i % 12)
		year := start.Year + (start.Month+i-1)/12

		for i := 1; i <= DaysPerMonth(month, year); i++ {
			dateRange = append(dateRange, &Date{i, month, year})
		}
	}

	// Add final month.
	for i := 1; i <= end.Day; i++ {
		dateRange = append(dateRange, &Date{i, end.Month, end.Year})
	}

	return dateRange, nil
}

func DaysPerMonth(month, year int) int {
	if month == 2 {
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			return 31
		}
	}
	return daysPerMonth[month-1]
}

func CurrentDate() *Date {
	t := time.Now()
	return &Date{t.Day(), int(t.Month()), t.Year()}
}
