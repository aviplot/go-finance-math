package financial

import (
	"github.com/aviplot/go-finance-math/test/testdata"
	"testing"
)

func TestPv(t *testing.T) {
	// Read known data.
	tdTab := testdata.TESTGetPvTestData()

	for _, td := range tdTab {
		// (rate float64, nper int64, pmt float64, fv float64, t bool)
		expected := td.Result
		precision := getPrecisionFromFloat(expected)

		result := Pv(td.Rate, td.Nper, td.Pmt, td.Fv, td.Type)
		result = round(result, precision)
		if result != expected {
			t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
		}
	}
}
