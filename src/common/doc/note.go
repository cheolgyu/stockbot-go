package doc

import (
	"github.com/cheolgyu/stockbot/src/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateNoteOne(field string, value string) (*mongo.UpdateResult, error) {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	coll := client.Database(DB_DATA).Collection(DB_DATA_COLL_PRICE)
	filter := bson.D{{}}
	update := bson.D{{"$set", bson.D{{field, value}}}}

	result, err := coll.UpdateOne(ctx, filter, update)

	return result, err
}
