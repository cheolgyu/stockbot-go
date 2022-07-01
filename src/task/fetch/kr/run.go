package kr

import (
	"github.com/cheolgyu/stockbot/src/fetch/kr/company"
	"github.com/cheolgyu/stockbot/src/fetch/kr/price"
)

func Run() {
	cr := company.Run{}
	cr.Run()
	pr := price.Run{}
	pr.Run()
}
