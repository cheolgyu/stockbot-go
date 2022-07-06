package vol

import (
	"context"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
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

	panic(" datainfo 항목 재구성으로 가격콜렉션 리셋필요.    ")

	//companys := doc.GetCompany()
	// for _, v := range companys[:1] {
	// 	prices := select_new_price_data(v.Code.Code, v.AggVolAt)

	// 	years := make(map[int]bool)
	// 	for _, v := range prices {
	// 		years[v.DateInfo.Dt_y] = true

	// 	}
	// 	insert_new_price_data(prices)
	// }
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

// 2. 코드별 마지막가격 데이터의 년도에 해당하는 가격데이터의 조회한다.
// func select_price_list (code string) []model.PriceMarket {
// 	filter := bson.M{"code": code, "dt": bson.M{"$lt": last}}
// 	opts := options.Find().SetSort(bson.M{"dt": -1})

// 	var results []model.PriceMarket
// 	cursor, err := collection_price.Find(context.TODO(), filter, opts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err = cursor.All(context.TODO(), &results); err != nil {
// 		panic(err)
// 	}
// 	return results

// }
