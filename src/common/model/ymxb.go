package model

import (
	"math"
)

type Unit_quote struct {
	Exchange
	Tick  int
	Price int `bson:"price"`
}

// 소수점2자리, 반올림
func Round2(x float64) float64 {
	return math.Round(x*100) / 100
}

// 소수점0자리, 반올림
func Round0(x float64) int {
	return int(math.Round(x))
}

type YmxbType int

const (
	//p1은 마지막 바운스점, p2는 마지막 가격점
	YmxbType1 YmxbType = iota
	//저가의 바운스점중 p1은 뒤에서 2번째 점, p2는 뒤에서 1번째 점
	YmxbType2
	//고가의 바운스점중 p1은 뒤에서 2번째 바운스점, p2는 뒤에서 1번째 바운스점
	YmxbType3
)

/*
	y=mx+b,
	p1은 저고종시가의 마지막 바운스점,
	p2는 저고종시가의 마지막 점.
	p3는 다음점

	p1~3는 실제 가격이고 tp1~3은 추상화한 가격,
	tp의 x는 일자를 1칸으로 추상화하였고,
	tp의 y는 호가를 1칸으로 추상화하였음.

	+ 현재는 p1은 마지막 바운스고 p2는 마지막점을 type1로 두고
	+ 나중에 type2는 마지막 저점 바운스 2개고
	+ type3는 마지막 고점 바운스 2개 로 만들기.
*/
type Ymxb struct {
	Code      string
	YmxbType  `bson:"ymxb_type"`
	PriceType `bson:"price_type"`

	P1 Point
	P2 Point
	P3 Point

	M float64
	B float64
}
