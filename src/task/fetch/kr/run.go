package kr

import (
	"github.com/cheolgyu/stockbot/src/fetch/kr/company"
)

func Run() {
	cr := company.Run{}
	cr.Run()
	// pr := price.Run{}
	// pr.Run()
}
