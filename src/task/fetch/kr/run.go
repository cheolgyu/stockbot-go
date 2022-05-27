package kr

import (
	"github.com/cheolgyu/stockbot/src/fetch/kr/company"
)

func Run() {
	cr := company.Run{}
	cr.Exe()
	// kr_price()
}

// func kr_price() {

// 	//list := price.SelectCodeAll()
// 	sd, ed := price.StartEndDate()
// 	codes := price.SelectCodeAll()
// 	price.NaverChart{}

// }
