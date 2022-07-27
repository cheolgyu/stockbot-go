package main

import (
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/task/asmb/line/bound"
	"github.com/cheolgyu/stockbot/src/task/asmb/line/ymxb"
)

func main() {
	line_bound := bound.LineBound{}
	line_bound.BaseRunStart(mlog.LineBound)
	line_bound.EXE()
	line_bound.BaseRunEnd()

	line_ymxb := ymxb.LineYmxb{}
	line_ymxb.BaseRunStart(mlog.LineYmxb)
	line_ymxb.EXE()
	line_ymxb.BaseRunEnd()
}
