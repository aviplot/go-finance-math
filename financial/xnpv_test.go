package financial

import (
	"github.com/aviplot/go-finance-math/test/testdata"
	"testing"
)

// TestXnpv tests Xnpv
func TestXnpv(t *testing.T) {
	// Read known data.
	tdTab := testdata.TESTGetCashflowTestData()

	for _, td := range tdTab {
		// Create cashflow using test data.
		cf := NewCashFlowTab(td.Amount, td.DateStart, td.IncomeTimes, td.Income, td.DateIncomeStart)
		//fmt.Println(cf)

		expected := td.ExpectedNPV
		precision := getPrecisionFromFloat(expected)

		result, _ := Xnpv(td.Rate, cf)
		result = round(result, precision)
		if result != expected {
			t.Fatalf("Error, result: \"%v\" Expected: \"%v\"", result, expected)
		}
	}
}
