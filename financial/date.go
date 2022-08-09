package financial

import (
	"encoding/json"
	"time"
)

// layout is how date is represented as string.
const layout = "2006-01-02"

// Date represent Date, neglecting time.
type Date struct {
	Date time.Time // Using standard time
}

// NewDateFromTime returns Date from standard time (removes the time out if the Date)
func NewDateFromTime(t time.Time) (d Date) {
	d.Date = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return
}

// NewDateFromFormattedString converts "yyyy-mm-dd" to Date
func NewDateFromFormattedString(ts string) (d Date) {
	t, err := time.Parse(layout, ts)
	if err != nil {
		// Error converting the date
		panic(err)
	}
	d.Date = t
	return

	//s := strings.Split(ts, "-")
	//yyyy, _ := strconv.Atoi(s[0])
	//mm, _ := strconv.Atoi(s[1])
	//dd, _ := strconv.Atoi(s[2])
	//
	//d.Date = time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	//return
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

func (d Date) String() string {
	//return fmt.Sprintf("%04d-%02d-%02d", d.Date.Year(), int(d.Date.Month()), d.Date.Day())
	return d.Date.Format(layout)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
