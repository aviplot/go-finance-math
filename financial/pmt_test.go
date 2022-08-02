package financial

import (
	"testing"
)

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
