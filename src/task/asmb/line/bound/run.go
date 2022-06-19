package bound

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
			bline := BoundLine{
				PriceType: v2,
				Code:      v.Code,
			}

			bline.GetStartingPoint()
			bline.GetAfterStartingPointPipeline()
			log.Println("bline.StartingPoint", bline.startingPoint)
			log.Println("bline.AfterStartingPoint ", bline.afterStartingPoint[:3])
			log.Println("bline.AfterStartingPoint len", len(bline.afterStartingPoint))
			bline.SetBoundPoint()
			log.Println("bline.SetBoundPoint len", len(bline.boundPoint))
			log.Println("bline.SetBoundPoint len", bline.boundPoint[:3])
			//bound point 찾기
		}
	}
}
