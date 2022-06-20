package model

import (
	"errors"
	"strings"
)

type Market int

const (

	//코스피
	KOSPI Market = iota
	//코스닥
	KOSDAQ Market = iota
	//코넥스
	KONEX Market = iota
)

var Market_arr = [...]Market{
	KOSPI,
	KOSDAQ,
	KONEX,
}
var Market_String = [...]string{
	"KOSPI",
	"KOSDAQ",
	"KONEX",
}

func String2Market(str string) (Market, error) {
	up := strings.ToUpper(str)
	ii := -10
	for i, v := range Market_String {
		if v == up {
			ii = i

		}
	}
	if ii >= 0 {
		return Market_arr[ii], nil
	}
	return Market_arr[0], errors.New("알수없는 마켓문자열입니다. " + str)
}

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
