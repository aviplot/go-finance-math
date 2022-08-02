package financial

import (
	"github.com/aviplot/go-finance-math/test/testdata"
	"testing"
)

func TestCoefficient(t *testing.T) {
	// Read known data.
	tdTab := testdata.TESTGetCashflowTestData()

	for _, td := range tdTab {
		// Create cashflow using test data.
		/*
			first float64, fDate string, months int, in float64, inDate string
		*/
		//fmt.Println(cf)
		expected := td.Coefficient
		precision := getPrecisionFromFloat(expected)
		t.Logf("Precision: %v", precision)
		result := ReturnCoefficient(td.Rate, int64(td.IncomeTimes), td.Amount, 0, false)
		t.Logf("result: %v", result)
		result = round(result, precision)
		if result != expected {
			t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
		}
	}
}
