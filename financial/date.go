package financial

import (
	"strconv"
	"strings"
	"time"
)

// Date represent Date, neglecting time.

type Date struct {
	Date time.Time // Using standard time
}

// NewDateFromTime returns Date from standard time (removes the time out if the Date)
func NewDateFromTime(t time.Time) (d Date) {
	d.Date = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return
}

// NewDateFromFormattedString converts "dd-mm-yyyy" to Date
func NewDateFromFormattedString(ts string) (d Date) {
	s := strings.Split(ts, "-")
	dd, _ := strconv.Atoi(s[0])
	mm, _ := strconv.Atoi(s[1])
	yyyy, _ := strconv.Atoi(s[2])
	d.Date = time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	return
}

// daysBetweenDates calculate amount of days between dates, neglecting time.
func daysBetweenDates(d1, d2 Date) int64 {
	h := d2.Date.Sub(d1.Date).Hours()
	days := h / 24
	return int64(days)
}

func (d Date) DaysFrom(d1 Date) int64 {
	return daysBetweenDates(d1, d)
}
func (d Date) DaysTo(d1 Date) int64 {
	return daysBetweenDates(d, d1)
}

func (d Date) AddMonth() (dr Date) {
	dr.Date = d.Date.AddDate(0, 1, 0)
	return
}
