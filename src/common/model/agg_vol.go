package model

type PeriodType int

const (
	//주
	Weeks PeriodType = iota
	//월
	Month PeriodType = iota
	//분기
	Quarter PeriodType = iota
)

type AggVolSum struct {
	Code       string
	Year       int
	VolWeeks   int
	VolMonth   int
	VolQuarter int
}
