package bp

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
)

type Run struct {
	code []model.Code
}

func (o *Run) Run() {
	o.code = doc.GetCodes()
	for _, v := range o.code[:1] {

		for _, v2 := range model.PriceTypes_arr {
			log.Println(v, v2)
			bl := BoundLine{
				PriceType: v2,
				Code:      v.Code,
			}

			bl.GetLastPoint()
			bl.GetAfterPoint()
			//bound point 찾기
		}
	}
}
