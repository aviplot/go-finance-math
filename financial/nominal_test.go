package financial

import (
	"testing"
)

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
