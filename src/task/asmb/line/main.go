package main

import (
	"log"

	"github.com/cheolgyu/stockbot/src/task/asmb/line/bp"
)

func main() {
	log.Println("i am line")
	bp_run := bp.Run{}
	bp_run.Run()

}
