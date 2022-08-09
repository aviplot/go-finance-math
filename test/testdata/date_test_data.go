package testdata

type dateData struct {
	Date          string
	DatePlusMonth string
	TargetDate    string
	DaysToTarget  int64
}

func TESTGetDateData() []dateData {
	return []dateData{
		{
			Date:          "1982-05-19",
			DatePlusMonth: "1982-06-19",
			TargetDate:    "2020-02-20",
			DaysToTarget:  13791,
		},
		{
			Date:          "2022-08-31",
			DatePlusMonth: "2022-10-01", // Note that we passed two months
			TargetDate:    "2022-09-01",
			DaysToTarget:  1,
		},
	}
}
