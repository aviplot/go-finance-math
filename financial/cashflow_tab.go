package financial

import "errors"

type cashFlowTab []cashFlow

var ErrEmptySlice = errors.New("slice is empty")

func NewCashFlowTab(first float64, fDate string, months int, in float64, inDate string) (result cashFlowTab) {
	f := cashFlow{
		NewDateFromFormattedString(fDate), first,
	}
	result = append(result, f)

	currentMonth := NewDateFromFormattedString(inDate)
	for months > 0 {
		result = append(result, cashFlow{
			currentMonth, in,
		})
		currentMonth = currentMonth.AddMonth()
		months--
	}
	return
}

func (ca cashFlowTab) FirstFlow() float64 {
	if len(ca) > 0 {
		return ca[0].Flow
	}
	panic(ErrEmptySlice)
}

func (ca cashFlowTab) FirstDate() date {
	if len(ca) > 0 {
		return ca[0].Date
	}
	panic(ErrEmptySlice)
}

func (ca cashFlowTab) String() (result string) {
	for _, c := range ca {
		result = result + c.String() + "\n"
	}
	return
}
