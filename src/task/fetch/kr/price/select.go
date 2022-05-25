package price

import (
	"fmt"
	"log"
	"time"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/kr/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartEndDate() (start string, end string) {

	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	coll := client.Database(config.DB_PUB).Collection(config.DB_PUB_COLL_NOTE)
	opts := options.Find().SetProjection(bson.M{"kr_price_updated_date": 1})
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)
	var data []bson.M
	cursor.All(ctx, &data)
	if len(data) == 0 {
		log.Println("kr_price_updated_date 신규")
		result, err := coll.InsertOne(ctx, bson.M{"kr_price_updated_date": config.PRICE_DEFAULT_START_DATE})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("kr_price_updated_date 신규:저장후 id  ", result.InsertedID)
		fmt.Println("1.start=", config.PRICE_DEFAULT_START_DATE)
	} else {

		start = data[0]["kr_price_updated_date"].(string)
		log.Println("kr_price_updated_date 존재하는 값 ", start)
	}

	end = time.Now().Format(config.PRICE_DATE_FORMAT)
	fmt.Println("end=", end)

	return start, end
}

func SelectCodeAll() []model.Company {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	cursor, err := client.Database(config.DB_PUB).Collection(config.DB_PUB_COLL_COMPANY).Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	if err != nil {
		log.Fatal(err)
	}
	var list []model.Company

	cursor.All(ctx, &list)

	fmt.Println("select company_list len=", len(list))
	return list
}
