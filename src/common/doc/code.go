package doc

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCompany() []model.Company {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	cursor, err := client.Database(DB_PUB).Collection(DB_PUB_COLL_COMPANY).Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Company

	defer cursor.Close(ctx)
	cursor.All(ctx, &list)

	return list
}

func GetCompanyCodes() []model.Company {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	opts := options.Find().SetProjection(bson.M{"code": 1, "name": 1, "update_price": 1})

	cursor, err := client.Database(DB_PUB).Collection(DB_PUB_COLL_COMPANY).Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Company

	defer cursor.Close(ctx)
	cursor.All(ctx, &list)

	return list
}
