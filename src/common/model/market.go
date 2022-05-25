package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

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

func parseUint(str string) (int, error) {
	// 08 일경우 오류 발생.
	res, err := strconv.Atoi(str)
	return int(res), err
}

func (o *PriceMarket) StringToPrice(str string) {

	arr := strings.Split(str, ",")
	var s0 = arr[0]

	if res, err := parseUint(s0); err == nil {
		o.Dt = res
	} else if err != nil {
		panic(err)
	}

	if res, err := parseUint(s0[:4]); err == nil {
		o.Dt_y = res
	} else if err != nil {
		panic(err)
	}

	if res, err := parseUint(s0[4:6]); err == nil {
		o.Dt_m = res
	} else if err != nil {
		log.Println("err:", str)
		log.Println("err:", str)
		panic(err)
	}

	if res, err := parseUint(s0[4:6]); err == nil {
		if res, err := convert_g4(res); err == nil {
			o.Dt_q4 = res
		} else if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[1], 32); err == nil {
		o.OpenPrice = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[2], 32); err == nil {
		o.HighPrice = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[3], 32); err == nil {
		o.LowPrice = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[4], 32); err == nil {
		o.ClosePrice = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := parseUint(arr[5]); err == nil {
		o.Volume = res
	} else if err != nil {
		panic(err)
	}

	str_fr := strings.Replace(arr[6], ",", "", -1)
	if str_fr == "" {
		str_fr = "0"
	}
	o.ForeignerBurnoutRate = str_fr

}
