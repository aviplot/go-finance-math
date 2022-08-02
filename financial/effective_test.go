package financial

import (
	"testing"
)

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
