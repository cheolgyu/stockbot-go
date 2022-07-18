package main

import (
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/task/asmb/agg/vol"
)

func main() {
	mlog.Info(mlog.AggVol, "main start")
	vol.Run()
}
