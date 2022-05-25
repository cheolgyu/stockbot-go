package price

import (
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Run struct {
	companys []model.Company
	start    string
	end      string
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 company의 마지막가격일자 갱신
*/
func (o *Run) Run() {
	o.getCodeCompany()

	start, end := StartEndDate()
	o.start = start
	o.end = end

}

func (o *Run) getCodeCompany() {
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
