package financial

import (
	"strconv"
	"strings"
	"time"
)

// Date represent date, neglecting time.

type date struct {
	Date time.Time // Using standard time
}

// NewDateFromTime returns date from standard time (removes the time out if the date)
func NewDateFromTime(t time.Time) (d date) {
	d.Date = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return
}

// NewDateFromFormattedString converts "dd-mm-yyyy" to date
func NewDateFromFormattedString(ts string) (d date) {
	s := strings.Split(ts, "-")
	dd, _ := strconv.Atoi(s[0])
	mm, _ := strconv.Atoi(s[1])
	yyyy, _ := strconv.Atoi(s[2])
	d.Date = time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	return
}

// daysBetweenDates calculate amount of day between dates, neglecting time.
func daysBetweenDates(d1, d2 date) int64 {
	h := d2.Date.Sub(d1.Date).Hours()
	days := h / 24
	return int64(days)
}

func (d date) DaysFrom(d1 date) int64 {
	return daysBetweenDates(d1, d)
}
func (d date) DaysTo(d1 date) int64 {
	return daysBetweenDates(d, d1)
}

func (d date) AddMonth() (dr date) {
	dr.Date = d.Date.AddDate(0, 1, 0)
	return
}
