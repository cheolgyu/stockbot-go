package doc

import (
	"context"
	"log"

	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCompany(client *mongo.Client, country model.Country) []model.Company {
	filter := bson.M{}
	if country != model.ALL {
		filter = bson.M{
			"country": country,
		}
	}

	cursor, err := client.Database(DB_PUB).Collection(DB_PUB_COLL_COMPANY).Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Company

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &list)

	return list
}

func GetCodes(client *mongo.Client, country model.Country) []model.Code {
	opts := options.Find().SetProjection(bson.M{"code": 1, "name": 1})
	// filter := bson.M{"country": country}
	filter := bson.M{}
	if country != model.ALL {
		filter = bson.M{
			"country": country,
		}
	}

	cursor, err := client.Database(DB_PUB).Collection(DB_PUB_COLL_COMPANY).Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Code

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &list)

	return list
}

func GetCodesMarket(client *mongo.Client, country model.Country) []model.Code {
	opts := options.Find().SetProjection(bson.M{"code": 1, "name": 1})
	filter := bson.M{}
	if country != model.ALL {
		filter = bson.M{
			"country": country,
		}
	}

	cursor, err := client.Database(DB_PUB).Collection(DB_PUB_COLL_MARKET).Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Code

	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &list)

	return list
}
