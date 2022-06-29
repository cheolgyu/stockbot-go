package price

import (
	"context"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Run struct {
	code     []model.Code
	start    string
	end      string
	download naverChart
	_        Insert
	openings []interface{}
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 pub.note 마지막가격일자 갱신
*/
func (o *Run) Run() {

	create_index()
	create_index_opening()
	o.code = doc.GetCodes()
	// 마켓 데이터필요
	//o.code = append(o.code, add_market()...)
	o.start, o.end = startEndDate()

	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	//dataPriceCollection := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
	open := make(map[int]int)
	for _, v := range o.code[:10] {

		o.download = naverChart{
			startDate: o.start,
			endDate:   o.end,
			Code:      v,
		}
		list, err := o.download.Run()
		if err != nil {
			log.Panic(err)
		}
		for k, _ := range o.download.Openings {
			open[k] = k
		}

		if len(list) > 0 {
			// insert := Insert{
			// 	coll: dataPriceCollection,
			// 	code: v,
			// 	list: list,
			// }
			//insert.Run()

		} else {
			log.Println(" price data size 0", v)
		}

	}

	for k, _ := range open {
		o.openings = append(o.openings, model.NewOpening(model.KR, k))
	}
	isnert_opening(o.openings)

	_, err := doc.UpdateNoteOne(doc.DB_PUB_COLL_NOTE_PRICE_UPDATED_KR, o.end)
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

func add_market() []model.Code {
	kr_market := []model.Code{
		{
			Code: "KOSPI",
			Name: "코스피",
		}, {
			Code: "KOSDAQ",
			Name: "코스닥",
		},
	}
	return kr_market
}
