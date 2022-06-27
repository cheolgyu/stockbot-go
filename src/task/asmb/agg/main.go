package main

import (
	"log"

	"github.com/cheolgyu/stockbot/src/task/asmb/agg/vol"
)

func main() {
	log.Println("i am agg")
	vol.Run()
}
