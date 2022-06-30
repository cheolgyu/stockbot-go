package bound

import (
	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func init() {
	client, _ = common.Connect()
}

type Run struct {
	code []model.Code
}

func (o *Run) Run() {
	o.code = doc.GetCodes(client, model.KR)
	for _, v := range o.code {

		for _, v2 := range model.PriceTypes_arr {
			bline := BoundLine{
				PriceType: v2,
				Code:      v.Code,
			}

			bline.GetStartingPoint()
			bline.GetAfterStartingPointPipeline()
			bline.SetBoundPoint()
			if len(bline.boundPoint) > 0 {
				bline.Insert()
			}

		}
	}
}
