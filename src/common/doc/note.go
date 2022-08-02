package doc

import (
	"time"

	"github.com/cheolgyu/stockbot/src/common"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateNoteOne(field string) {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	value := time.Now().Format("2006-01-02 15:04:05")
	coll := client.Database(DB_PUB).Collection(DB_PUB_COLL_NOTE)
	filter := bson.D{{}}
	update := bson.D{{"$set", bson.D{{field, value}}}}

	_, err := coll.UpdateOne(ctx, filter, update)
	panic(err)

}
