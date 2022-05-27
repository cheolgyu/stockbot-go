package company

import (
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/kr/config"
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
	companyCollection := client.Database(config.DB_PUB).Collection(config.DB_PUB_COLL_COMPANY)
	//var ui []interface{}
	models := []mongo.WriteModel{}

	for _, v := range o.Company {
		models = append(models, mongo.NewReplaceOneModel().SetFilter(bson.M{"_id": v.ID, "code": v.Code}).SetUpsert(true).SetReplacement(v))
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
