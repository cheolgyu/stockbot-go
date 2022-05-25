package kr

import (
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
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
	companyCollection := client.Database("test").Collection("company")
	//var ui []interface{}
	models := []mongo.WriteModel{}

	for _, v := range o.Company {
		log.Println(v.ID)
		//ui = append(ui, v)
		var vi interface{}
		vi = v

		models = append(models, mongo.NewReplaceOneModel().SetFilter(bson.M{"_id": v.ID, "code": v.Code}).SetUpsert(true).SetReplacement(vi))
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := companyCollection.BulkWrite(ctx, models, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"inserted %v and deleted %v documents\n",
		res.InsertedCount,
		res.DeletedCount)
	// result, err := companyCollection.InsertMany(ctx, ui)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
	// for _, id := range result.InsertedIDs {
	// 	fmt.Printf("\t%s\n", id)
	// }
}
