package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Market struct {
	Code Code `bson:"inline"`
	Country
	UpdatedMarketAt int `bson:"updated_market_at"`
}

var Loc *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	Loc = loc
}

type Exchange int

const (
	/*
		한국
	*/
	//코스피
	KOSPI Exchange = iota
	//코스닥
	KOSDAQ
	//코넥스
	KONEX
	/*
		미국
	*/
	//나스닥
	NASDAQ
	//뉴욕증권거래소
	NYSE
	//아멕스(아메리칸 익스프레스)
	AMEX
)

var Exchanges = map[Country]map[Exchange]Code{
	KR: {
		KOSPI:  {"KOSPI", "코스피"},
		KOSDAQ: {"KOSDAQ", "코스닥"},
		KONEX:  {"KONEX", "코넥스"},
	},
	US: {
		NASDAQ: {"NASDAQ", "나스닥"},
		NYSE:   {"NYSE", "뉴욕증권거래소"},
		AMEX:   {"AMEX", "아멕스"},
	}}

func ConvertExchanges(country Country, code string) (Exchange, error) {
	upcode := strings.ToUpper(code)
	var key Exchange
	for k, v := range Exchanges[country] {
		if v.Code == upcode {
			return k, nil
		}
	}
	return key, errors.New("알수없는 마켓문자열입니다. " + upcode)
}

type Price struct {
	Display string
	Decimal uint
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
	OP       primitive.Decimal128
	CP       primitive.Decimal128
	LP       primitive.Decimal128
	HP       primitive.Decimal128
	Vol      int
	//ForeignerBurnoutRate
	FBR string
}

// "3.87" to 3.87
func ParsePrice(myString string) primitive.Decimal128 {
	p, err := primitive.ParseDecimal128(myString)
	if err != nil {
		panic(err)
	}
	return p
}

func ParseVol(myString string) int {

	s := strings.Replace(myString, ",", "", -1)
	value, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return value

}

type Country string

const (
	//전체
	ALL Country = "all"
	//한국
	KR Country = "kr"
	//미국
	US Country = "us"
)

var Countrys = map[string]Country{
	"kr": KR,
	"us": US,
}

type DateInfo struct {
	Dt      int
	Y       int
	M       int
	D       int
	Week    int
	Quarter int
}

// dt: 20220725
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
		msg := fmt.Sprintf(" 분기 변환 오류: %d", num)
		err = errors.New(msg)
	}

	return res, err
}
