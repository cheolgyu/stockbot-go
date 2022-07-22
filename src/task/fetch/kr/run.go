package kr

import (
	"github.com/cheolgyu/stockbot/src/fetch/kr/price"
)

func Run() {

	pr := price.Run{Downlad: true}
	pr.Run()
}
