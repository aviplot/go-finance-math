package testdata

type CashFlowTestData struct {
	Rate            float64
	Amount          float64
	DateStart       string // "dd-mm-yyyy"
	Income          float64
	DateIncomeStart string // "dd-mm-yyyy"
	IncomeTimes     int
	Balloon         float64
	BalloonDate     string // "dd-mm-yyyy"
	ExpectedIRR     float64
	ExpectedNPV     float64
	Pv              float64
	Coefficient     float64
}

func TESTGetCashflowTestData() []CashFlowTestData {
	return []CashFlowTestData{
		{
			Rate:            0.04, // 4%
			Amount:          -100000,
			DateStart:       "15-05-2000",
			Income:          4000,
			DateIncomeStart: "20-06-2000",
			IncomeTimes:     36,
			Balloon:         0,
			BalloonDate:     "20-06-2000",
			ExpectedIRR:     0.2826335,
			ExpectedNPV:     35543.5544,
			Pv:              -130993.51,
			Coefficient:     -33.870766422,
		},
		{
			Rate:            0.03, // 3%
			Amount:          -100000,
			DateStart:       "15-05-2000",
			Income:          6000,
			DateIncomeStart: "20-06-2000",
			IncomeTimes:     36,
			Balloon:         0,
			BalloonDate:     "20-06-2000",
			ExpectedIRR:     0.769936717,
			ExpectedNPV:     106347.1928,
			Pv:              111,
			Coefficient:     -34.3864651,
		},
		{
			Rate:            0.15, // 3%
			Amount:          -20000,
			DateStart:       "15-05-2000",
			Income:          1247.6957453033,
			DateIncomeStart: "20-06-2000",
			IncomeTimes:     18,
			Balloon:         0,
			BalloonDate:     "20-06-2000",
			ExpectedIRR:     0.15719098,
			ExpectedNPV:     97.61421509,
			Pv:              0,
			Coefficient:     -16.029548931,
		},
	}
}
