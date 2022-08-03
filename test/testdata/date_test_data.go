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
			Date:          "19-05-1982",
			DatePlusMonth: "19-06-1982",
			TargetDate:    "20-02-2020",
			DaysToTarget:  13791,
		},
		{
			Date:          "31-08-2022",
			DatePlusMonth: "01-10-2022", // Note that we passed two months
			TargetDate:    "01-09-2022",
			DaysToTarget:  1,
		},
	}
}
