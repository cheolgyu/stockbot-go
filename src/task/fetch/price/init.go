package price

import (
	"context"
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common/doc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const INDEX_NAME_OPENING string = "index_country_1_dt_1"
const INDEX_NAME_PRICE string = "index_code_1_dt_1"

//https://christiangiacomi.com/posts/mongodb-index-using-go/
func create_index_price(client *mongo.Client) {
	log.Println("create_index_price() start")
	dataPriceCollection := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)

	exist := false
	specs, err := dataPriceCollection.Indexes().ListSpecifications(context.TODO())
	for _, v := range specs {
		if v.Name == INDEX_NAME_PRICE {
			exist = true
		}
	}
	if err != nil {
		panic(err)
	}
	if !exist {
		mod := mongo.IndexModel{
			Keys:    bson.D{{"code", 1}, {"dt", 1}}, // index in ascending order or -1 for descending order
			Options: options.Index().SetUnique(true).SetName(INDEX_NAME_PRICE),
		}

		_, err := dataPriceCollection.Indexes().CreateOne(context.TODO(), mod)
		if err != nil {
			// 5. Something went wrong, we log it and return false
			fmt.Println(err.Error())
			panic(err.Error())
		}
	}
	log.Println("create_index_price() end")
}

func create_index_opening(client *mongo.Client) {
	log.Println("create_index_opening() start")
	coll := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE_OPENING)

	exist := false
	specs, err := coll.Indexes().ListSpecifications(context.TODO())
	for _, v := range specs {
		if v.Name == INDEX_NAME_OPENING {
			exist = true
		}
	}
	if err != nil {
		panic(err)
	}

	if !exist {
		mod := mongo.IndexModel{
			Keys:    bson.D{{"country", 1}, {"dt", 1}}, // index in ascending order or -1 for descending order
			Options: options.Index().SetUnique(true).SetName(INDEX_NAME_OPENING),
		}
		_, err := coll.Indexes().CreateOne(context.TODO(), mod)
		if err != nil {
			// 5. Something went wrong, we log it and return false
			fmt.Println(err.Error())
			panic(err.Error())
		}
	}
	log.Println("create_index_opening() end")
}
