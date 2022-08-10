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

// TestPmt validate PMT function
func TestPmt(t *testing.T) {
	expected := -580.54
	precision := getPrecisionFromFloat(expected)
	result := Pmt(0.05, 36, 10000, 500, true) // rate , nper, pv , fv , t
	result = round(result, precision)
	if result != expected {
		t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
	}
}

// TestNominal validate GetNominal function
func TestNominal(t *testing.T) {
	expected := 0.058411
	precision := getPrecisionFromFloat(expected)
	result, _ := GetNominal(0.06, 12)
	result = round(result, precision)
	if result != expected {
		t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
	}
}

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

// TestNominal validate GetNominal function
func TestEffective(t *testing.T) {
	expected := 0.061677812
	precision := getPrecisionFromFloat(expected)
	result, _ := GetEffectiveRate(0.06, 12)
	result = round(result, precision)
	if result != expected {
		t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
	}
}

// TestNominal validate GetNominal function
func TestCashFlow(t *testing.T) {
	r := 0.12345
	cf := NewCashFlowPayments(1000000, "2010-05-10", 240, r/12)
	result, e := Xirr(cf)
	if e != nil {
		t.Fatalf("Error: %v", e)
		return
	}
	expected := 0.13053038213
	precision := getPrecisionFromFloat(expected)
	result = round(result, precision)
	expected = round(expected, precision)
	if result != expected {
		t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
	}
}

// TestIrr validate IRR function
func TestIrr(t *testing.T) {
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
		result, _ := Irr(cf)
		result = round(result, precision)
		expected =  round(expected, precision)
		if result != expected {
			t.Fatalf("Error, result: \"%v\" Expected: \"%v\" (Precision: %v)", result, expected, precision)
		}
	}
}