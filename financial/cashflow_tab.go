package financial

import (
	"errors"
	"sort"
)

type CashFlowTab []CashFlow

var ErrEmptySlice = errors.New("slice is empty")

// NewCashFlowTab created new CashFlowTab, used in testcases.
func NewCashFlowTab(first float64, fDate string, months int, in float64, inDate string) (result CashFlowTab) {
	f := CashFlow{
		NewDateFromFormattedString(fDate), first,
	}
	result = append(result, f)

	currentMonth := NewDateFromFormattedString(inDate)
	for months > 0 {
		result = append(result, CashFlow{
			currentMonth, in,
		})
		currentMonth = currentMonth.AddMonth()
		months--
	}
	return
}

func (ca CashFlowTab) FirstFlow() float64 {
	if len(ca) > 0 {
		return ca[0].Flow
	}
	panic(ErrEmptySlice)
}

func (ca CashFlowTab) FirstDate() Date {
	if len(ca) > 0 {
		return ca[0].Date
	}
	panic(ErrEmptySlice)
}

func (ca CashFlowTab) String() (result string) {
	for _, c := range ca {
		result = result + c.String() + "\n"
	}
	return
}

// Len impl "Interface" to support sorting, using sort.Sort.
func (ca CashFlowTab) Len() int {
	return len(ca)
}

// Swap impl "Interface" to support sorting, using sort.Sort.
func (ca CashFlowTab) Swap(i, j int) {
	ca[i], ca[j] = ca[j], ca[i]
}

// Less impl "Interface" to support sorting, using sort.Sort.
func (ca CashFlowTab) Less(i, j int) bool {
	return ca[i].Date.Date.Before(ca[j].Date.Date)
}

func (ca CashFlowTab) OrderByDate() (r CashFlowTab) {
	r = ca
	sort.Sort(r)
	return
}
