package testdata

type pvTestData struct {
	Rate   float64
	Nper   int64
	Pmt    float64
	Fv     float64
	Type   bool
	Result float64
}

func TESTGetPvTestData() []pvTestData {
	return []pvTestData{
		{
			Rate:   0.05,
			Nper:   36,
			Pmt:    -6544.345,
			Fv:     1212.44,
			Type:   true,
			Result: 113493.38,
		},
		{
			Rate:   0.15,
			Nper:   17,
			Pmt:    -2133.345,
			Fv:     1212.44,
			Type:   true,
			Result: 14723.12,
		},
		{
			Rate:   0.15,
			Nper:   17,
			Pmt:    -2133.345,
			Fv:     1212.44,
			Type:   false,
			Result: 12788.01,
		},
	}
}
