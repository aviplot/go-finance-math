package financial

import (
	"fmt"
	"math"
	"strings"
)

// round is rounding a number with another int param which is the number of digits after the dot.
func round(f float64, r int) float64 {
	p := math.Pow(10, float64(r))
	return math.Round(f*p) / p
}

// getPrecisionFromFloat analyze float and get the precision.
func getPrecisionFromFloat(f float64) int {
	_, frac := math.Modf(math.Abs(f))
	s := fmt.Sprintf("%.10f", frac)
	s = strings.Split(s, ".")[1] // Taking the fraction
	for {
		if strings.HasSuffix(s, "0") {
			s = s[:len(s)-1]
		} else {
			break
		}
	}

	return len(s)
}
