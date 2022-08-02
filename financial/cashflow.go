package financial

import (
	"fmt"
)

type cashFlow struct {
	Date date
	Flow float64
}

type cashFlowWithCalculated struct {
	cashFlow
	OrderId   int
	Period    int
	Interest  float64
	Principal float64
}

func NewCashFlow() cashFlow {
	return cashFlow{}
}

func NewCashFlowWithCalculated() cashFlowWithCalculated {
	return cashFlowWithCalculated{}
}

func (c cashFlow) String() string {
	return fmt.Sprintf("Date: %v | flow: %v", c.Date, c.Flow)
}
