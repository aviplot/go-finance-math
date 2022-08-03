package testdata

type fvTestData struct {
	Rate   float64
	Nper   int64
	Pmt    float64
	Pv     float64
	Type   bool
	Result float64
}

func TESTGetFvTestData() []fvTestData {
	return []fvTestData{
		{
			Rate:   0.05,
			Nper:   36,
			Pmt:    -6544.345,
			Pv:     1212.44,
			Type:   true,
			Result: 651523.03,
		},
		{
			Rate:   0.15,
			Nper:   17,
			Pmt:    -2133.345,
			Pv:     1212.44,
			Type:   true,
			Result: 146604.38,
		},
		{
			Rate:   0.15,
			Nper:   17,
			Pmt:    -2133.345,
			Pv:     1212.44,
			Type:   false,
			Result: 125780.24,
		},
	}
}
