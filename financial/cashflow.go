package financial

import (
	"fmt"
)

// cashFlow Holds one record of a cash flow table.
type cashFlow struct {
	Date date
	Flow float64
}

// NewCashFlow returns new empty cashflow record
func NewCashFlow() cashFlow {
	return cashFlow{}
}

func (c cashFlow) String() string {
	return fmt.Sprintf("Date: %v | flow: %v", c.Date, c.Flow)
}
