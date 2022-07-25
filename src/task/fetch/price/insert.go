package price

import (
	"context"
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Insert struct {
	coll *mongo.Collection
	code model.Code
	list []model.PriceMarket
}

func (o *Insert) Run() {
	//var ui []interface{}
	models := []mongo.WriteModel{}

	for _, v := range o.list {
		models = append(models, mongo.NewReplaceOneModel().SetFilter(bson.M{"code": o.code.Code, "dt": v.DateInfo.Dt}).SetUpsert(true).SetReplacement(v))
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := o.coll.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"%+v: inserted %v documents and upserted %v documents and ModifiedCount %v and deleted %v documents\n",
		o.code,
		res.InsertedCount,
		res.UpsertedCount,
		res.ModifiedCount,
		res.DeletedCount)

}
