package main

import (
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/task/asmb/agg/vol"
)

func main() {

	item := vol.AggVol{}
	item.BaseRunStart(mlog.AggVol)
	item.EXE()
	item.BaseRunEnd()
}
