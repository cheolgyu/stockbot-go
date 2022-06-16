package bound

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BoundLine struct {
	PriceType          model.PriceType
	Code               string
	startingPoint      model.Point
	afterStartingPoint []model.Point
	boundPoint         []model.Bound
}

func (o *BoundLine) GetStartingPoint() {

	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	projection := bson.M{"p2_x": 1, "p2_y": 1}
	filter := bson.M{"code": o.Code, "p": o.PriceType}
	opts := options.Find().SetSort(bson.M{"p2_x": 1}).SetLimit(1).SetProjection(projection)

	cursor, err := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_BOUND_POINT).Find(ctx, filter, opts)
	if err != nil {
		log.Panicln(err.Error())
	}

	var list []model.Point
	if err = cursor.All(ctx, &list); err != nil {
		log.Panicln(err.Error())
	}
	defer cursor.Close(ctx)

	if len(list) == 0 {
		o.startingPoint = model.Point{0, 0}
	} else {
		o.startingPoint = list[0]
	}
}

func (o *BoundLine) GetAfterStartingPointPipeline() {
	var res []model.Point
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	matchStage := bson.D{{"$match", bson.D{{"code", o.Code}, {"dt", bson.D{{"$gt", o.startingPoint.X}}}}}}
	projectStage := bson.D{{"$project", bson.D{{"_id", 0}, {"x", "$dt"}, {"y", "$" + o.PriceType.String_DB_Field()}}}}
	sortStage := bson.D{{"$sort", bson.D{{"dt", 1}}}}

	cursor, err := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE).Aggregate(ctx, mongo.Pipeline{matchStage, projectStage, sortStage})
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &res); err != nil {
		log.Panicln(err.Error())
	}

	o.afterStartingPoint = res
}

func (o *BoundLine) SetBoundPoint() {
	if len(o.afterStartingPoint) > 2 {
		return
	}

	j := o.startingPoint
	d := 0

	for i := 0; i < len(o.afterStartingPoint); i++ {
		v := o.afterStartingPoint[i]

		if j.Y < v.Y {
			d = model.Increase.Int()
		} else if j.Y > v.Y {
			d = model.Decrease.Int()
		} else {
			d = model.Constant.Int()
		}

	}

}
