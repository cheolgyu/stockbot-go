package price

import (
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Run struct {
	code     []model.Code
	start    string
	end      string
	download naverChart
	_        Insert
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 pub.note 마지막가격일자 갱신
*/
func (o *Run) Run() {
	create_index()
	o.code = doc.GetCodes()
	o.code = append(o.code, add_market()...)
	o.start, o.end = startEndDate()

	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	dataPriceCollection := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)

	for _, v := range o.code {
		o.download = naverChart{
			startDate: o.start,
			endDate:   o.end,
			Code:      v,
		}
		list, err := o.download.Run()
		if err != nil {
			log.Panic(err)
		}
		if len(list) > 0 {
			insert := Insert{
				coll: dataPriceCollection,
				code: v,
				list: list,
			}
			insert.Run()

		} else {
			log.Println(" price data size 0", v)
		}

	}

	res, err := doc.UpdateNoteOne(doc.DB_DATA_COLL_PRICE_UPDATED_KR, o.end)
	if err != nil {
		panic(err.Error())
	}
	log.Println("doc.UpdateNoteOne:", res)
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

//https://christiangiacomi.com/posts/mongodb-index-using-go/
func create_index() {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	dataPriceCollection := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)

	/*
			// IndexModel represents a new index to be created.
		type IndexModel struct {
			// A document describing which keys should be used for the index. It cannot be nil. This must be an order-preserving
			// type such as bson.D. Map types such as bson.M are not valid. See https://docs.mongodb.com/manual/indexes/#indexes
			// for examples of valid documents.
			Keys interface{}

			// The options to use to create the index.
			Options *options.IndexOptions
		}
	*/
	mod := mongo.IndexModel{
		Keys:    bson.D{{"code", 1}, {"dt", 1}}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := dataPriceCollection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		panic(err.Error())
	}
}
