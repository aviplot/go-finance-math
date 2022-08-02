package financial

import (
	"github.com/aviplot/go-finance-math/test/testdata"
	"testing"
)

// TestXirr validate XIRR function
func TestXirr(t *testing.T) {
	// Read known data.
	tdTab := testdata.TESTGetCashflowTestData()

	for _, td := range tdTab {
		// Create cashflow using test data.
		/*
			first float64, fDate string, months int, in float64, inDate string
		*/
		cf := NewCashFlowTab(td.Amount, td.DateStart, td.IncomeTimes, td.Income, td.DateIncomeStart)
		//fmt.Println(cf)
		expected := td.ExpectedIRR
		precision := getPrecisionFromFloat(expected)
		result, _ := Xirr(cf)
		result = round(result, precision)
		if result != expected {
			t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
		}
	}
}
