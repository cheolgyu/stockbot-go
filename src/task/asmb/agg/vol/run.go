package vol

import (
	"context"
	"math"
	"time"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/schollz/progressbar/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection_price *mongo.Collection
var collection_agg_vol *mongo.Collection
var collection_agg_vol_sum *mongo.Collection

/*
절차
	agg_vol_sum 계산범위:  코드별 마지막 가격데이터의 년도
	1. agg_vol_sum에서 테이블의 코드별 마지막 년도를 구한다 없으면0을 반환한다.
		1.1 agg_vol_sum에서 테이블의 코드별 마지막 년도를 구한다 없으면0을 반환한다.
		1.2 agg_vol_sum의 year가 0이면 code의 price 데이터중 가장 작은 year은 구한다.
		1.3 특정연도의 가격테이터를 조회한다.
	2. 코드별 연도별 가격데이터에서 주별 거래량의합, 월별 거래량의 합, 분기별 거래량의합을 구한다.
	3. 코드별 해당연도의 agg_vol_sum을 upsert 한다.
	4. 코드별 전체연도의 agg_vol_sum을 조회한다.
	5. 코드별 코드의 전체연도의 agg_vol_sum데이터로 기간별 표준편차를 구한다.
		5.1 표준편차를 구한다.
			5.1.1 관찰값들의 평균을 구한다. (편차값을 구하기 위해서)
			5.1.2 편차를 구한다. (편차: 관측값에서 평균을 뺀것)
			5.1.3 표준편차를 구한다.
		5.2 빈도를 분석한다.
	6. 코드별 코드의 전체연도의 기간별 표준편차를 저장한다.

*/
func Run() {
	var ctx context.Context
	client, ctx = common.Connect()
	defer client.Disconnect(ctx)

	collection_price = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
	collection_agg_vol = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_AGG_VOL)
	collection_agg_vol_sum = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_AGG_VOL_SUM)

	t := time.Now()
	cur_year := t.Year()
	companys := doc.GetCompany(client)

	bar := progressbar.Default(int64(len(companys)))

	for _, v := range companys {
		bar.Add(1)
		year_start := 0
		for year_item := 0; year_item <= cur_year; year_item++ {

			if year_item == 0 {
				year_start = chk_year_agg_vol_sum(v.Code.Code)
				if year_start == 0 {
					year_item = chk_year_price(v.Code.Code)
				} else {
					year_item = year_start
				}
			}

			prices := select_price_eq_year(v.Code.Code, year_item)
			if len(prices) != 0 {
				agg_vol_sum := sum(prices)
				agg_vol_sum.Calculate()
				upsert_sum(agg_vol_sum)
			}

		}

		total_sum := select_total_agg_vol_sum(v.Code.Code)
		aggVol := get_percent(total_sum)
		upsert_agg_vol(aggVol)
	}
}

// 1.1 agg_vol_sum에서 테이블의 코드별 마지막 년도를 구한다 없으면0을 반환한다.
func chk_year_agg_vol_sum(code string) (year int) {
	year = -1
	noDocument := false

	matchStage := bson.D{{"$match", bson.D{{"code", code}}}}
	groupStage := bson.D{
		{"$group",
			bson.D{
				{"_id", ""},
				{"max_y", bson.D{{"$max", "$y"}}},
			},
		},
	}

	pipeline := mongo.Pipeline{matchStage, groupStage}

	cursor, err := collection_price.Aggregate(context.TODO(), pipeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			year = 0
			noDocument = true
		} else {
			mlog.Err(mlog.AggVol, err)
			panic(err)
		}
	}

	if !noDocument {
		var results []model.AggVolSum
		if err = cursor.All(context.TODO(), &results); err != nil {
			mlog.Err(mlog.AggVol, err)
			panic(err)
		}
		year = results[0].Year
	}

	return year

}

//1-2 agg_vol_sum의 year가 0이면 code의 price 데이터중 가장 작은 year은 구한다.
func chk_year_price(code string) (year int) {

	filter := bson.M{"code": code}
	sort := bson.M{"dt": 1}

	opts := options.FindOne()
	opts.SetSort(sort)

	res := model.PriceMarket{}
	collection_price.FindOne(context.TODO(), filter, opts).Decode(&res)

	year = res.Y

	return year
}

//1.3 특정연도의 가격테이터를 조회한다.
func select_price_eq_year(code string, year int) (prices []model.PriceMarket) {
	matchStage := bson.D{{"$match", bson.D{{"code", code}, {"y", year}}}}
	sortStage := bson.D{
		{"$sort",
			bson.D{
				{"dt", 1},
			},
		},
	}

	pipeline := mongo.Pipeline{matchStage, sortStage}

	cursor, err := collection_price.Aggregate(context.TODO(), pipeline)
	if err != nil {
		mlog.Err(mlog.AggVol, err)
		panic(err)
	}

	if err = cursor.All(context.TODO(), &prices); err != nil {
		mlog.Err(mlog.AggVol, err)
		panic(err)
	}

	return prices
}

//2. 코드별 가격데이터에서 주별 거래량의합, 월별 거래량의 합, 분기별 거래량의합을 구한다.
func sum(list []model.PriceMarket) model.AggVolSum {
	item := model.AggVolSum{}
	item.Code = list[0].Code
	item.Year = list[0].Y
	var w map[int]int = make(map[int]int)
	var m map[int]int = make(map[int]int)
	var q map[int]int = make(map[int]int)

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

// 3. 코드별 해당연도의 agg_vol_sum을 upsert 한다.
func upsert_sum(v model.AggVolSum) {
	filter := bson.M{"code": v.Code, "year": v.Year}
	opts := options.ReplaceOptions{}
	opts.SetUpsert(true)

	_, error := collection_agg_vol_sum.ReplaceOne(context.TODO(), filter, v, &opts)
	if error != nil {
		panic(error)
	}
}

//4. 코드별 전체연도의 agg_vol_sum을 조회한다.
func select_total_agg_vol_sum(code string) []model.AggVolSum {
	var res []model.AggVolSum

	filter := bson.M{"code": code}
	sort := bson.M{"year": 1}
	opts := options.Find()
	opts.SetSort(sort)
	cur, err := collection_agg_vol_sum.Find(context.TODO(), filter, opts)
	if err != nil {
		mlog.Err(mlog.AggVol, err)
		panic(err)
	}

	err = cur.All(context.TODO(), &res)
	if err != nil {
		mlog.Err(mlog.AggVol, err)
		panic(err)
	}

	return res
}

//5.1 표준편차를 구한다.
//(방법: https://namu.wiki/w/표준편차 )
func get_percent(list []model.AggVolSum) model.AggVol {
	var res model.AggVol
	res.Result = make(map[string]model.Statistic)
	res.Code = list[0].Code

	for _, ov := range model.ObservationValueTypeArr {
		avgObsrvValue, data := get_avg_ObservationValue(ov, list)
		deviation_list := get_deviation(ov, avgObsrvValue, list)
		std_dvtn, square_avg := get_standard_deviation(deviation_list)

		avs := model.Statistic{}
		avs.BaseAnalysis.DataCnt = len(deviation_list)
		avs.BaseAnalysis.Data = data
		avs.BaseAnalysis.Average = avgObsrvValue
		avs.BaseAnalysis.Variance = square_avg
		avs.BaseAnalysis.StandardDeviation = std_dvtn

		fa := model.FrequencyAnalysis{}
		fa.Table = getFrequencyAnalysis(ov, data)
		avs.FrequencyAnalysis = fa

		res.Result[ov.ToField()] = avs
	}

	return res
}

//5.2 빈도를 분석한다.
func getFrequencyAnalysis(ov model.ObservationValueType, data map[int]int) map[int]model.FrequencyTableRow {
	item := make(map[int]model.FrequencyTableRow)

	all_frequency := 0
	for _, v := range data {
		row := item[v]
		row.Value = v
		row.Frequency++
		item[v] = row
		all_frequency++
	}

	for k, _ := range item {
		row := item[k]
		row.Percentage = model.Round2((float64(row.Frequency) / float64(all_frequency)) * 100)
		item[k] = row
	}

	return item
}

// 5.1.1 관찰값들의 평균을 구한다. (편차값을 구하기 위해서)
func get_avg_ObservationValue(ov model.ObservationValueType, list []model.AggVolSum) (avg int, data map[int]int) {
	data = make(map[int]int)

	sum := 0
	cnt := 0
	for _, v := range list {

		obsrvValue := v.GetValueOfObservationValueType(ov)
		data[v.Year] = obsrvValue

		sum += obsrvValue
		cnt++
	}
	if cnt != 0 {
		avg = sum / cnt
	}

	return avg, data
}

// 6-2 편차를 구한다. (편차: 관측값에서 평균을 뺀것)
func get_deviation(ov model.ObservationValueType, avg int, list []model.AggVolSum) (res []int) {
	for _, v := range list {

		obsrvValue := v.GetValueOfObservationValueType(ov)
		deviation := avg - obsrvValue
		res = append(res, deviation)
	}
	return res
}

/*
6-3 표준편차를 구한다.
  편차들의 '제곱'을 구한 뒤 그 편차들의 제곱의 평균을 구해서 나온 값[5]에
  다시 제곱근을 구하는 우회적인 방법을 써서 산포도 값을 구한 것이 바로 표준편차
*/
func get_standard_deviation(deviation_list []int) (standard_deviation float64, square_avg float64) {
	var square_list []int

	square_sum := 0
	square_cnt := 0
	for _, v := range deviation_list {
		vv := v * v
		square_sum += vv
		square_list = append(square_list, vv)
		square_cnt++
	}
	square_avg = float64(square_sum) / float64(square_cnt)
	standard_deviation = math.Sqrt(square_avg)

	return standard_deviation, square_avg
}

//6. 코드별 코드의 전체연도의 기간별 표준편차를 저장한다.
func upsert_agg_vol(item model.AggVol) {

	filter := bson.M{"code": item.Code}
	opts := options.ReplaceOptions{}
	opts.SetUpsert(true)

	_, error := collection_agg_vol.ReplaceOne(context.TODO(), filter, item, &opts)
	if error != nil {
		mlog.Err(mlog.AggVol, error)
		panic(error)
	}
}
