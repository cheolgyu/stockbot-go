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

	var p1 model.Point = o.startingPoint
	if p1.X == 0 {
		p1 = o.afterStartingPoint[0]
		o.afterStartingPoint = o.afterStartingPoint[1:]
	}

	count_p2 := len(o.afterStartingPoint)
	if count_p2 < 2 {
		return
	}

	var prev_direction model.Direction

	var duration uint

	for i := 0; i < count_p2; i++ {

		var cur_direction model.Direction

		p2 := o.afterStartingPoint[i]

		if p1.Y < p2.Y {
			cur_direction = model.Increase
		} else if p1.Y > p2.Y {
			cur_direction = model.Decrease
		} else {
			cur_direction = model.Constant
		}

		chage_direction := prev_direction != cur_direction
		last_p2 := i+1 == count_p2

		if i != 0 && (chage_direction || last_p2) {

			bound := model.Bound{
				Direction: cur_direction,
				Duration:  duration,
				P1:        p1,
				P2:        p2,
			}
			o.boundPoint = append(o.boundPoint, bound)
			p1 = p2
			prev_direction = cur_direction
			duration = 0
		} else {
			p1 = p2
			prev_direction = cur_direction
		}

		duration++

	}

}
