package vol

import (
	"context"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var collection_price *mongo.Collection
var collection_agg_vol *mongo.Collection
var collection_agg_vol_sum *mongo.Collection

/*
절차
	agg_vol_sum 계산범위:  코드별 마지막 가격데이터의 년도
	1. 코드별 마지막가격 데이터의 년도에 해당하는 가격데이터의 조회한다.
	2. 코드별 가격데이터에서 주별 거래량의합, 월별 거래량의 합, 분기별 거래량의합을 구한다.
	3. 코드별 해당연도의 agg_vol_sum을 upsert 한다.
	4. 코드별 전체연도의 agg_vol_sum을 조회한다.
	5. 코드별 코드의 전체연도의 agg_vol_sum데이터로 기간별 퍼센트 값목록을 구한다.
	6. 코드별 코드의 전체연도의 기간별 퍼센트 값목록을 저장한다.

*/
func Run() {
	var ctx context.Context
	client, ctx = common.Connect()
	defer client.Disconnect(ctx)

	collection_price = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
	collection_agg_vol = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_AGG_VOL)
	collection_agg_vol_sum = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_AGG_VOL_SUM)

	companys := doc.GetCompany(client)
	for _, v := range companys[:10] {
		prices := select_price_list(v.Code.Code)
		agg_vol_sum := sum(prices)
		log.Println(v, agg_vol_sum)
		//panic("가격 데이터 dt 작은순으로 출력하면 2022년도가 제일먼저나옴;;;; 가격데이터가 잘 안들어간듯")
	}
}

// 1. 코드별 마지막가격 데이터의 년도에 해당하는 가격데이터의 조회한다.
func select_price_list(code string) []model.PriceMarket {

	matchStage := bson.D{{"$match", bson.D{{"code", code}}}}
	groupStage := bson.D{
		{"$group",
			bson.D{
				{"_id", ""},
				{"max_y", bson.D{{"$max", "$y"}}},
			},
		},
	}
	lookupStage := bson.D{
		{"$lookup",
			bson.D{
				{"from", "price"},
				{"localField", "max_y"},
				{"foreignField", "y"},
				{"pipeline",
					bson.A{
						bson.D{{"$match", bson.D{{"code", code}}}},
						bson.D{
							{"$project",
								bson.D{
									{"_id", 0},
									{"cp", 0},
									{"op", 0},
									{"lp", 0},
									{"hp", 0},
									{"fbr", 0},
								},
							},
						},
						bson.D{{"$sort", bson.D{{"dt", -1}}}},
					},
				},
				{"as", "list"},
			},
		},
	}

	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$list"}}}}
	replaceRootStage := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$list"}}}}

	pipeline := mongo.Pipeline{matchStage, groupStage, lookupStage, unwindStage, replaceRootStage}

	var results []model.PriceMarket
	cursor, err := collection_price.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results

}

//2. 코드별 가격데이터에서 주별 거래량의합, 월별 거래량의 합, 분기별 거래량의합을 구한다.
func sum(list []model.PriceMarket) model.AggVolSum {
	item := model.AggVolSum{}

	var w []int = make([]int, 54)
	var m []int = make([]int, 13)
	var q []int = make([]int, 5)

	for _, v := range list {
		w[v.DateInfo.Week] += v.Vol
		m[v.DateInfo.M] += v.Vol
		q[v.DateInfo.Quarter] += v.Vol
	}
	item.SumWeeks = w
	item.SumMonth = m
	item.SumQuarter = q

	return item
}

// func insert_new_price_data(price []model.PriceMarket) {

// 	models := []mongo.WriteModel{}

// 	for _, v := range price {
// 		i := model.AggVolSum{Code: v.Code,Year: v.Dt_y,PeriodType: ,VOL}
// 		filter := bson.M{"code": i.Code, "price_type": i.PriceType, "p1.x": i.P1.X}
// 		models = append(models, mongo.NewReplaceOneModel().SetFilter(filter).SetUpsert(true).SetReplacement(i))
// 	}

// 	opts := options.BulkWrite().SetOrdered(false)
// 	res, err := collection_agg_vol_sum.BulkWrite(context.TODO(), models, opts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf(
// 		"inserted %v documents and upserted %v documents and ModifiedCount %v and deleted %v documents\n",
// 		res.InsertedCount,
// 		res.UpsertedCount,
// 		res.ModifiedCount,
// 		res.DeletedCount)

// }
