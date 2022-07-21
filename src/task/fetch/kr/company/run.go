package company

import (
	"context"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/dao"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
기존 목록 조회
new 목록 다운 및 변환
기존 목록에 new 목록으로 바꾸기
replace방식으로 upsert설정 후 저장
*/
type Run struct {
	old_list map[string]model.Company
	downlad  Req_krx
	convert  Convert
}

var client *mongo.Client

func (o *Run) Run() {
	client, _ = common.Connect()
	defer client.Disconnect(context.TODO())

	o.old_list = dao.SelectMapCompany(client, model.KR)
	log.Println("len(old_companys)=", len(o.old_list))

	o.downlad = Req_krx{}
	o.downlad.Run()

	o.convert.Old = o.old_list
	o.convert.Run()

	var list []model.Company
	for _, v := range o.old_list {
		list = append(list, o.old_list[v.Code.Code])
	}

	log.Println("len(merge_companys)=", len(o.old_list))
	dao.Insert(client, list)
}
