package financial

import (
	"github.com/aviplot/go-finance-math/test/testdata"
	"testing"
)

func BenchmarkXirr(b *testing.B) {
	// Read known data.
	tdTab := testdata.TESTGetCashflowTestData()
	td := tdTab[0]
	cf := NewCashFlowTab(td.Amount, td.DateStart, td.IncomeTimes, td.Income, td.DateIncomeStart)

	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		_, err := Xirr(cf)
		if err != nil {
			b.Fatalf("Xirr failed")
		}
	}
}

func BenchmarkCashflowCalculation(b *testing.B) {
	// Read known data.
	tdTab := testdata.TESTGetCashflowTestData()
	var cftm cashFlowTabMulti
	for _, td := range tdTab {
		cft := NewCashFlowTab(td.Amount, td.DateStart, td.IncomeTimes, td.Income, td.DateIncomeStart)
		cftm = append(cftm, cft)
	}

	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		calculated := NewCalculatedCashFlow(cftm)
		for _, ccc := range calculated {
			b.Logf("irr: %v", ccc.calc.Xirr)
		}
		//if  {
		//	b.Fatalf("Xirr failed")
		//}
	}

}
