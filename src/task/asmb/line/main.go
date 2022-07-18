package main

import (
	"log"

	"github.com/cheolgyu/stockbot/src/task/asmb/line/bound"
	"github.com/cheolgyu/stockbot/src/task/asmb/line/ymxb"
)

func main() {
	log.Println("i am line")
	bp_run := bound.Run{}
	bp_run.Run()
	ymxb.Run()

}
