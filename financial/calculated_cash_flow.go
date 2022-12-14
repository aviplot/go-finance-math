package financial

// calculatedCashFlow is cashflow, with added calculated information
type calculatedCashFlow struct {
	CashFlow
	OrderId   int
	Period    int
	Interest  float64
	Principal float64
}

func NewEmptyCalculatedCashFlow() calculatedCashFlow {
	return calculatedCashFlow{}
}
