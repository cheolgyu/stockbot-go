package price

import (
	"context"
	"fmt"

	"github.com/cheolgyu/stockbot/src/common/doc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//https://christiangiacomi.com/posts/mongodb-index-using-go/
func create_index(client *mongo.Client) {
	dataPriceCollection := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)

	/*
			// IndexModel represents a new index to be created.
		type IndexModel struct {
			// A document describing which keys should be used for the index. It cannot be nil. This must be an order-preserving
			// type such as bson.D. Map types such as bson.M are not valid. See https://docs.mongodb.com/manual/indexes/#indexes
			// for examples of valid documents.
			Keys interface{}

			// The options to use to create the index.
			Options *options.IndexOptions
		}
	*/
	mod := mongo.IndexModel{
		Keys:    bson.D{{"code", 1}, {"dt", 1}}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := dataPriceCollection.Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		panic(err.Error())
	}
}

func create_index_opening(client *mongo.Client) {

	coll := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE_OPENING)

	mod := mongo.IndexModel{
		Keys:    bson.D{{"country", 1}, {"dt", 1}}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := coll.Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		panic(err.Error())
	}

}
