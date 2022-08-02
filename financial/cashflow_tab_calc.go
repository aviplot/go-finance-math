package financial

type calculatedData struct {
	Xirr float64
	//Pv float64
}

type cashflowTabCalc struct {
	cfTab cashFlowTab
	calc  calculatedData
}

type cashflowTabCalcTab []cashflowTabCalc

func NewCalculatedCashFlow(cftm cashFlowTabMulti) (result cashflowTabCalcTab) {

	for _, cft := range cftm {
		x, _ := Xirr(cft)
		record := cashflowTabCalc{
			cfTab: cft,
			calc: calculatedData{
				Xirr: x,
			},
		}
		result = append(result, record)
	}
	return
}
