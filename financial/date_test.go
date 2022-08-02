package financial

import (
	"github.com/aviplot/go-finance-math/test/testdata"
	"log"
	"testing"
)

// TestDate validate date format YYYY-MM-DD
func TestAddMonth(t *testing.T) {
	// Read known data.
	dt := testdata.TESTGetDateData()

	for _, d := range dt {
		baseDate := NewDateFromFormattedString(d.Date)
		monthTestDate := NewDateFromFormattedString(d.DatePlusMonth)
		targetDate := NewDateFromFormattedString(d.TargetDate)

		basePMDate := baseDate.AddMonth()
		if !basePMDate.Date.Equal(monthTestDate.Date) {
			log.Fatalf("Base date: %v +month is: %v but expected is: %v", baseDate, basePMDate, monthTestDate)
		}

		if baseDate.DaysTo(targetDate) != d.DaysToTarget {
			log.Fatalf("Days expected: %v but we got: %v", d.DaysToTarget, baseDate.DaysTo(targetDate))
		}
	}
}
