package price

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	kr_price "github.com/cheolgyu/stockbot/src/fetch/kr/price"
	us_price "github.com/cheolgyu/stockbot/src/fetch/us/price"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var collectionPrice *mongo.Collection

//custom := now.Format("2006-01-02 15:04:05")
const PRICE_DATE_FORMAT = "20060102"

func init() {
	client, ctx = common.Connect()
	create_index_price(client)
	create_index_opening(client)
	collectionPrice = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
}

type CrawlingPrice interface {
	GetResult(downlad bool) ([]model.PriceMarket, error)
}

//key:code, value:last_price_date
var startDate map[string]string
var endDate string

type RunCode struct {
	code     model.Code
	download CrawlingPrice
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 pub.note 마지막가격일자 갱신
*/
func Run() {
	startDate, endDate = startEnd()

	for _, v := range model.Countrys {
		run_country(v)
	}

	defer client.Disconnect(ctx)

}

func run_country(country model.Country) {

	code := doc.GetCodes(client, country)
	code = append(code, doc.GetCodesMarket(client, country)...)

	for _, v := range code {
		run_code := RunCode{code: v}

		switch country {
		case model.KR:
			run_code.download = &kr_price.NaverChart{
				StartDate: startDate[run_code.code.Code],
				EndDate:   endDate,
				Code:      v,
			}
		case model.US:
			run_code.download = &us_price.RequestNasdaqCom{
				StartDate: startDate[run_code.code.Code],
				Code:      v,
			}
		}

		list, err := run_code.download.GetResult(false)
		if err != nil {
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

			log.Println(" price data size 0", v)
			fmt.Printf(" price data is 0 ,%+v \n", run_code)
		}

	}

	_, err := doc.UpdateNoteOne(doc.DB_PUB_COLL_NOTE_PRICE_UPDATED_KR, endDate)
	if err != nil {
		panic(err.Error())
	}
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
		log.Fatal(err)
	}
	result := []struct {
		Code string
		Dt   int
	}{}
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range result {
		startDate[v.Code] = strconv.Itoa(v.Dt)
	}

	defer cursor.Close(ctx)

	endDate := time.Now().Format(PRICE_DATE_FORMAT)
	return startDate, endDate
}
