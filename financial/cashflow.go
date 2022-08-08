package financial

import (
	"fmt"
)

// CashFlow Holds one record of a cash flow table.
type CashFlow struct {
	Date Date    `json:"date"`
	Flow float64 `json:"flow"`
}

// NewCashFlow returns new empty cashflow record
func NewCashFlow() CashFlow {
	return CashFlow{}
}

func (c CashFlow) String() string {
	return fmt.Sprintf("Date: %v | flow: %v", c.Date, c.Flow)
}
