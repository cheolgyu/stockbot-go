package model

var MarketList = []string{"KOSPI", "KOSDAQ", "FUT", "KPI200"}
var MarketListName = []string{"코스피", "코스닥", "선물", "코스피200"}

type PriceMarket struct {
	Dt                   int
	Dt_y                 int
	Dt_m                 int
	Dt_q4                int
	OpenPrice            float32
	HighPrice            float32
	LowPrice             float32
	ClosePrice           float32
	Volume               int
	ForeignerBurnoutRate string
}
