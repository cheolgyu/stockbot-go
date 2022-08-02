package price

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/base"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/common/model"
	kr_price "github.com/cheolgyu/stockbot/src/fetch/country/kr/price"
	us_price "github.com/cheolgyu/stockbot/src/fetch/country/us/price"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var collectionPrice *mongo.Collection

//custom := now.Format("2006-01-02 15:04:05")
const PRICE_DATE_FORMAT = "20060102"

func init() {

	mlog.Info(mlog.Fetch, "start price/price.go init")
	client, ctx = common.Connect()
	create_index_price(client)
	create_index_opening(client)
	collectionPrice = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
	mlog.Info(mlog.Fetch, "end price/price.go init")
}

type CrawlingPrice interface {
	GetResult(downlad bool) ([]model.PriceMarket, error)
}

type FetchPrice struct {
	base.Run

	//key:code, value:last_price_date
	startDate map[string]string
	endDate   string
}

type runCode struct {
	code     model.Code
	download CrawlingPrice
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 pub.note 마지막가격일자 갱신
*/
func (o *FetchPrice) EXE() {
	//	delete_us_prices(client)
	o.startDate, o.endDate = startEnd()

	o.run_country()

	defer client.Disconnect(ctx)

}

func (o *FetchPrice) run_country() {

	code := doc.GetCodes(client, o.Country)
	code = append(code, doc.GetCodesMarket(client, o.Country)...)

	for _, v := range code {
		run_code := runCode{code: v}

		switch o.Country {
		case model.KR:
			run_code.download = &kr_price.NaverChart{
				StartDate: o.startDate[run_code.code.Code],
				EndDate:   o.endDate,
				Code:      v,
			}
		case model.US:
			run_code.download = &us_price.RequestNasdaqCom{
				StartDate: o.startDate[run_code.code.Code],
				Code:      v,
			}
		}

		list, err := run_code.download.GetResult(o.Download)
		if err != nil {
			mlog.Err(mlog.Fetch, err)
			log.Panic(err)
		}

		if len(list) > 0 {
			insert := Insert{
				coll: collectionPrice,
				code: v,
				list: list,
			}
			insert.Run()
		} else {
			mlog.Info(mlog.Fetch, " price data is 0 ,%+v \n", run_code)
			log.Println(" price data size 0", v)
			fmt.Printf(" price data is 0 ,%+v \n", run_code)
		}

	}

	doc.UpdateNoteOne(doc.GetNoteField(o.Country, doc.FETCH_PRICE_UPDATE))
}

//return (key:code value:20220725, 20220725 )
func startEnd() (map[string]string, string) {
	startDate := make(map[string]string)

	projectStage := bson.D{{"$project", bson.D{{"code", "$_id"}, {"dt", 1}}}}
	groupState := bson.D{
		{"$group",
			bson.D{
				{"_id", "$code"},
				{"dt", bson.D{{"$max", "$dt"}}},
			},
		},
	}

	cursor, err := collectionPrice.Aggregate(ctx, mongo.Pipeline{groupState, projectStage})
	if err != nil {
		mlog.Err(mlog.Fetch, err)
		log.Fatal(err)
	}
	result := []struct {
		Code string
		Dt   int
	}{}
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		mlog.Err(mlog.Fetch, err)
		log.Fatal(err)
	}

	for _, v := range result {
		startDate[v.Code] = strconv.Itoa(v.Dt)
	}

	defer cursor.Close(ctx)

	endDate := time.Now().Format(PRICE_DATE_FORMAT)
	return startDate, endDate
}

func delete_us_prices(client *mongo.Client) {
	log.Println("delete_us_prices() start")
	coll := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)

	code := doc.GetCodes(client, model.US)
	for _, v := range code {
		_, err := coll.DeleteMany(context.TODO(), bson.M{"code": v.Code})
		if err != nil {
			panic(err)
		}
	}

	log.Println("delete_us_prices() end")
}
