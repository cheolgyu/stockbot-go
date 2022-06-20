package model

import "math"

type Unit_quote struct {
	Market
	Tick  int
	Price int `bson:"price"`
}

// 소수점2자리, 반올림
func Round2(x float64) float64 {
	return math.Round(x*100) / 100
}
