package bp

import (
	"context"
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BoundLine struct {
	PriceType  model.PriceType
	Code       string
	LastPoint  model.Point
	AfterPoint []model.Point
}

func (o *BoundLine) GetLastPoint() {

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
		o.LastPoint = model.Point{0, 0}
	} else {
		o.LastPoint = list[0]
	}

	log.Print("GetLastPoint", o.LastPoint)

}

func (o *BoundLine) GetAfterPoint() {
	var res []model.Point

	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	projection := bson.M{"_id": 0, "dt": 1, o.PriceType.String_DB_Field(): 1}
	filter := bson.M{"code": o.Code, "dt": bson.M{"$gte": o.LastPoint.X}}
	opts := options.Find().SetSort(bson.M{"dt": 1}).SetProjection(projection)

	log.Println("projection", projection)
	log.Println("filter", filter)

	cursor, err := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE).Find(ctx, filter, opts)
	if err != nil {
		log.Panicln(err.Error())
	}
	// if err = cursor.All(ctx, &res); err != nil {
	// 	log.Panicln(err.Error())
	// }

	// 커서 all로 받지 못함 이유 필드명이 다이나믹해서 받으려면 stcut model.pricemarket으로 받아야함
	// 받고나서 point로 파서함.
	// 방법2 agg pipline addfile 로 받으면됨. find와 agg 성능차이는 ?
	//https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/
	for cursor.Next(context.TODO()) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	o.AfterPoint = res
	log.Print(",o.AfterPoint len", len(res))
	//log.Println(res)
}
