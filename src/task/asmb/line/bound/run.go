package bound

import (
	"context"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/base"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var bound_collection *mongo.Collection
var price_collection *mongo.Collection

func init() {
	client, _ = common.Connect()
	bound_collection = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_BOUND_POINT)
	price_collection = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
}

type BoundRun struct {
	code []model.Code
	base.Run
}

func (o *BoundRun) BoundRun() {
	o.Start()

	defer client.Disconnect(context.TODO())

	o.code = doc.GetCodes(client, o.Country)
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

	o.End()
}
