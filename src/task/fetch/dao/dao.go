package dao

import (
	"context"
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SelectMapCompany(client *mongo.Client, country model.Country) map[string]model.Company {
	list := doc.GetCompany(client, country)
	company_map := make(map[string]model.Company)

	for _, v := range list {
		company_map[v.Code.Code] = v
	}

	return company_map
}

func Insert(client *mongo.Client, list []model.Company) {
	companyCollection := client.Database(doc.DB_PUB).Collection(doc.DB_PUB_COLL_COMPANY)
	models := []mongo.WriteModel{}

	for _, v := range list {
		models = append(models, mongo.NewReplaceOneModel().SetFilter(bson.M{"code": v.Code.Code}).SetUpsert(true).SetReplacement(v))
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := companyCollection.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Company InsertedCount %v and ModifiedCount %v and deleted %v documents\n",
		res.InsertedCount,
		res.ModifiedCount,
		res.DeletedCount)

}
