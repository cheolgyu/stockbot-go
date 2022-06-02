package bound

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CodeLine struct {
	PriceType model.PriceType
	Code      model.Code
	LastBound model.Bound
}

func (o *CodeLine) GetLastBound() model.Bound {
	var res model.Bound

	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	filter := bson.M{"code": o.Code, "p": o.PriceType}
	opts := options.Find().SetSort(bson.M{"p2_x": 1}).SetLimit(1)
	var err error
	if err = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_BOUND_POINT).FindOne(ctx, filter, opts).Decode(&res); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Company

	defer cursor.Close(ctx)
	cursor.One(ctx, &list)

	o.LastBound = res
}
