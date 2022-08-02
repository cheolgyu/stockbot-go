package main

import (
	"context"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/company"
	"github.com/cheolgyu/stockbot/src/fetch/price"
	"go.mongodb.org/mongo-driver/bson"
)

const DOWNLOAD_COMPANY bool = true
const DOWNLOAD_PRICE bool = true

func init() {
	mlog.Info(mlog.Fetch, "start main.go init")
	insert_Exchanges()
	mlog.Info(mlog.Fetch, "end main.go init")
}

func main() {
	company := company.FetchCompany{}
	company.Download = DOWNLOAD_COMPANY
	company.BaseRunStart(mlog.FetchCompany)
	company.EXE()
	company.BaseRunEnd()

	price := price.FetchPrice{}
	price.Download = DOWNLOAD_PRICE
	price.BaseRunStart(mlog.FetchPrice)
	price.EXE()
	price.BaseRunEnd()
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
