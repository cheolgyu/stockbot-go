package main

import (
	"log"

	"github.com/cheolgyu/stockbot/src/fetch/kr"
)

func main() {
	log.Println("i am fetch")
	kr.Run()
}
