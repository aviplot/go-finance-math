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
	ExpectedXIRR     float64
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
			DateStart:       "2000-05-15",
			Income:          4000,
			DateIncomeStart: "2000-06-20",
			IncomeTimes:     36,
			Balloon:         0,
			BalloonDate:     "2000-06-20",
			ExpectedXIRR:     0.2826335,
			ExpectedIRR:     0.02121114161818,
			ExpectedNPV:     35543.5544,
			Pv:              -130993.51,
			Coefficient:     -33.870766422,
		},
		{
			Rate:            0.03, // 3%
			Amount:          -100000,
			DateStart:       "2000-05-15",
			Income:          6000,
			DateIncomeStart: "2000-06-20",
			IncomeTimes:     36,
			Balloon:         0,
			BalloonDate:     "2000-06-20",
			ExpectedXIRR:     0.769936717,
			ExpectedIRR:     0.049439495,
			ExpectedNPV:     106347.1928,
			Pv:              111,
			Coefficient:     -34.3864651,
		},
		{
			Rate:            0.15, // 3%
			Amount:          -20000,
			DateStart:       "2000-05-15",
			Income:          1247.6957453033,
			DateIncomeStart: "2000-06-20",
			IncomeTimes:     18,
			Balloon:         0,
			BalloonDate:     "2000-06-20",
			ExpectedXIRR:     0.15719098,
			ExpectedIRR:     0.0125,
			ExpectedNPV:     97.61421509,
			Pv:              0,
			Coefficient:     -16.029548931,
		},
	}
}
