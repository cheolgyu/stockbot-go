package main

import (
	"context"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/us"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	insert_Exchanges()
}

func main() {
	log.Println("i am fetch")
	//kr.Run()
	us.Run()
}

func insert_Exchanges() {

	client, _ := common.Connect()
	coll := client.Database(doc.DB_PUB).Collection(doc.DB_PUB_COLL_MARKET)

	for k, v := range model.Exchanges {
		for kk, vv := range v {
			if kk != model.KONEX {
				filter := bson.M{"code": vv.Code, "country": k}
				cnt, err := coll.CountDocuments(context.TODO(), filter)
				if err != nil {
					panic(err)
				}
				if cnt == 0 {
					cmp := model.Market{
						Code: model.Code{
							Code: vv.Code,
							Name: vv.Name,
						},
						Country: k,
					}
					_, err := coll.InsertOne(context.TODO(), cmp)
					if err != nil {
						panic(err)
					}
				}
			}
		}

	}
	defer client.Disconnect(context.TODO())
}
