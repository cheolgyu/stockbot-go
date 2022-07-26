package main

import (
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/task/asmb/line/bound"
	"github.com/cheolgyu/stockbot/src/task/asmb/line/ymxb"
)

func main() {
	bp_run := bound.BoundRun{}
	bp_run.SetCountry()
	bp_run.Who = string(mlog.LineBound)
	bp_run.BoundRun()

	mlog.Info(mlog.LineBound, "end")
	mlog.Info(mlog.LineYmxb, "start")
	ymxb.Run()
	mlog.Info(mlog.LineYmxb, "end")

}
