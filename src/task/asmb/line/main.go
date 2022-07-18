package main

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/task/asmb/line/bound"
	"github.com/cheolgyu/stockbot/src/task/asmb/line/ymxb"
)

func main() {
	log.Println("i am line")
	mlog.Info(mlog.LineBound, "start")
	bp_run := bound.Run{}
	bp_run.Run()
	mlog.Info(mlog.LineBound, "end")
	mlog.Info(mlog.LineYmxb, "start")
	ymxb.Run()
	mlog.Info(mlog.LineYmxb, "end")

}
