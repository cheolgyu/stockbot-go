package price

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var START_DT_PRICE_DATE = 19560303

type Price struct {
	companys []model.Company
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 company의 마지막가격일자 갱신
*/
func (o *Price) Run() {
	o.getCodeCompany()

}

func (o *Price) getDataPrice() {

	var lastDt int

	//custom := now.Format("2006-01-02 15:04:05")
	format := time.Now().Format("20060102")
	log.Println(format)
	dt, err := strconv.Atoi(format)
	if err != nil {
		panic(err)
	}
	lastDt = dt

	fmt.Println("lastDt=", lastDt)
}

func (o *Price) getCodeCompany() {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	opts := options.Find().SetProjection(bson.M{"code": 1, "name": 1, "update_price": 1})

	cursor, err := client.Database("test").Collection("company").Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.Company

	defer cursor.Close(ctx)
	cursor.All(ctx, &list)

	fmt.Println("select company_list len=", len(list))

	o.companys = list
}
