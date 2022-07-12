package company

import (
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Insert struct {
	Company []model.Company
}

func (o *Insert) Run() {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	companyCollection := client.Database(doc.DB_PUB).Collection(doc.DB_PUB_COLL_COMPANY)
	//var ui []interface{}
	models := []mongo.WriteModel{}

	for _, v := range o.Company {
		models = append(models, mongo.NewReplaceOneModel().SetFilter(bson.M{"code": v.Code.Code}).SetUpsert(true).SetReplacement(v))
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := companyCollection.BulkWrite(ctx, models, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"ModifiedCount %v and deleted %v documents\n",
		res.ModifiedCount,
		res.DeletedCount)

}
