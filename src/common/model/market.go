package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Market struct {
	Code Code `bson:"inline"`
	Country
	UpdatedMarketAt int `bson:"updated_market_at"`
}
type MarketType int

var Loc *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	Loc = loc
}

const (

	//코스피
	KOSPI MarketType = iota
	//코스닥
	KOSDAQ MarketType = iota
	//코넥스
	KONEX MarketType = iota
)

// key:MarketType, value:Code
var MarketType_map = map[MarketType]Code{
	KOSPI:  Code{"KOSPI", "코스피"},
	KOSDAQ: Code{"KOSDAQ", "코스닥"},
	KONEX:  Code{"KONEX", "코넥스"},
}

var MarketType_arr = [...]MarketType{
	KOSPI,
	KOSDAQ,
	KONEX,
}
var MarketType_code_String = [...]string{
	"KOSPI",
	"KOSDAQ",
	"KONEX",
}
var MarketType_name_String = [...]string{
	"코스피",
	"코스닥",
	"코넥스",
}

func String2Market(str string) (MarketType, error) {
	up := strings.ToUpper(str)
	ii := -10
	for i, v := range MarketType_code_String {
		if v == up {
			ii = i

		}
	}
	if ii >= 0 {
		return MarketType_arr[ii], nil
	}
	return MarketType_arr[0], errors.New("알수없는 마켓문자열입니다. " + str)
}

/*
	OP    시가
	HP    고가
	LP    저가
	CP    종가
	Vol   거래량
*/
type PriceMarket struct {
	Code     string
	DateInfo `bson:"inline"`
	OP       float32
	CP       float32
	LP       float32
	HP       float32
	Vol      int
	//ForeignerBurnoutRate
	FBR string
}

type Country string

const (
	//한국
	KR Country = "kr"
	//미국
	US Country = "us"
)

type DateInfo struct {
	Dt      int
	Y       int
	M       int
	D       int
	Week    int
	Quarter int
}

func NewDateInfo(dt int) DateInfo {
	o := DateInfo{}
	o.Dt = dt
	sdt := strconv.Itoa(dt)
	if res, err := parseUint(sdt[:4]); err == nil {
		o.Y = res
	} else if err != nil {
		panic(err)
	}

	if res, err := parseUint(sdt[4:6]); err == nil {
		o.M = res
		if res2, err := convert_g4(o.M); err == nil {
			o.Quarter = res2
		} else if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	if res, err := parseUint(sdt[6:8]); err == nil {
		o.D = res
	} else if err != nil {
		panic(err)
	}

	//해당날짜가 몇주쨰인지 구한다.
	t := time.Date(o.Y, time.Month(o.M), o.D, 12, 12, 12, 12, Loc)
	_, w := t.ISOWeek()
	o.Week = w

	return o

}

type Opening struct {
	Country
	DateInfo `bson:"inline"`
}

func NewOpening(c Country, dt int) Opening {
	o := Opening{}
	o.Country = c
	o.DateInfo = NewDateInfo(dt)

	return o

}
func parseUint(str string) (int, error) {
	// 08 일경우 오류 발생.
	res, err := strconv.Atoi(str)
	return int(res), err
}
func convert_g4(num int) (int, error) {
	var err error
	res := 0
	// select   case  when mm < 4 then 1 when mm < 7 then 2  when mm < 10 then 3 else 4 end as q4
	if num < 4 {
		res = 1
	} else if num < 7 {
		res = 2
	} else if num < 10 {
		res = 3
	} else if num < 13 {
		res = 4
	} else {

		err = errors.New(fmt.Sprintf(" 분기 변환 오류: %v", num))
	}

	return res, err
}
