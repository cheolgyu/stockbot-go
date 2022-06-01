package model

var MarketList = []string{"KOSPI", "KOSDAQ", "FUT", "KPI200"}
var MarketListName = []string{"코스피", "코스닥", "선물", "코스피200"}

/*
	OP    시가
	HP    고가
	LP    저가
	CP    종가
	Vol   거래량
*/
type PriceMarket struct {
	Code  string
	Dt    int
	Dt_y  int
	Dt_m  int
	Dt_q4 int
	OP    float32
	CP    float32
	LP    float32
	HP    float32
	Vol   int
	//ForeignerBurnoutRate
	FBR string
}
