package main

import (
	"context"
	"flag"
	"strings"

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
	countryPtr := flag.String("country", "kr", "input country value")
	low_country := strings.ToLower(*countryPtr)
	var country model.Country
	for k, v := range model.Countrys {
		if k == low_country {
			country = v
		}
	}
	mlog.Info(mlog.Fetch, "start main:", countryPtr)
	mlog.Info(mlog.Fetch, "start company, DOWNLOAD_COMPANY=", DOWNLOAD_COMPANY, countryPtr)
	company.Run(country, DOWNLOAD_COMPANY)
	mlog.Info(mlog.Fetch, "end company", countryPtr)
	mlog.Info(mlog.Fetch, "start price, DOWNLOAD_PRICE=", DOWNLOAD_PRICE, countryPtr)
	price.Run(country, DOWNLOAD_PRICE)
	mlog.Info(mlog.Fetch, "end price", countryPtr)
	mlog.Info(mlog.Fetch, "end main", countryPtr)
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
