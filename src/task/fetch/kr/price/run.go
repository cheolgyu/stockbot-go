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
	"github.com/cheolgyu/stockbot/src/fetch/kr/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var collectionPrice *mongo.Collection

func init() {
	client, ctx = common.Connect()
	create_index_price(client)
	create_index_opening(client)
	collectionPrice = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
}

type Run struct {
	code     []model.Code
	_        Insert
	openings []interface{}
	//key:code, value:last_price_date
	startDate map[string]string
	endDate   string
}

type RunCode struct {
	code     model.Code
	download naverChart
	openings []interface{}
}

func (o *Run) setDate_StartEnd() {
	o.startDate = make(map[string]string)

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
		o.startDate[v.Code] = strconv.Itoa(v.Dt)
	}

	defer cursor.Close(ctx)

	o.endDate = time.Now().Format(config.PRICE_DATE_FORMAT)
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 pub.note 마지막가격일자 갱신
*/
func (o *Run) Run() {

	defer client.Disconnect(ctx)

	o.setDate_StartEnd()
	o.code = doc.GetCodes(client, model.KR)
	o.code = append(o.code, doc.GetCodesMarket(client, model.KR)...)

	open := make(map[int]int)
	for _, v := range o.code {
		run_code := RunCode{code: v}

		run_code.download = naverChart{
			startDate: o.startDate[run_code.code.Code],
			endDate:   o.endDate,
			Code:      v,
		}
		list, err := run_code.download.Run()
		if err != nil {
			log.Panic(err)
		}
		for k, _ := range run_code.download.Openings {
			open[k] = k
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
			fmt.Printf("%+v\n", run_code)
		}

	}

	for k, _ := range open {
		o.openings = append(o.openings, model.NewOpening(model.KR, k))
	}
	if len(o.openings) > 0 {
		isnert_opening(o.openings)
	}

	_, err := doc.UpdateNoteOne(doc.DB_PUB_COLL_NOTE_PRICE_UPDATED_KR, o.endDate)
	if err != nil {
		panic(err.Error())
	}

}

func isnert_opening(list []interface{}) {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	coll := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE_OPENING)
	_, err := coll.InsertMany(context.TODO(), list)
	if err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			panic(err.Error())
		}

	}
}
