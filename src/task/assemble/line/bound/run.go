package bound

import (
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
)

type Run struct {
	code []model.Code
}

func (o *Run) Run() {
	o.code = doc.GetCodes()

	for _, v := range o.code {
		v.Code
	}
}
