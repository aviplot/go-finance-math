package financial

import (
	"github.com/aviplot/go-finance-math/test/testdata"
	"testing"
)

func TestFv(t *testing.T) {
	// Read known data.
	tdTab := testdata.TESTGetFvTestData()

	for _, td := range tdTab {
		expected := td.Result
		precision := getPrecisionFromFloat(expected)

		result := Fv(td.Rate, td.Nper, td.Pmt, td.Pv, td.Type)
		result = round(result, precision)
		if result != expected {
			t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
		}
	}
}
