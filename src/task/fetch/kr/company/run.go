package company

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/model"
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
	insert   Insert
}

func (o *Run) Run() {
	o.old_list = company_map()
	log.Println("len(old_companys)=", len(o.old_list))

	o.downlad = Req_krx{}
	o.downlad.Run()

	o.convert.Old = o.old_list
	o.convert.Run()

	var list []model.Company
	for _, v := range o.old_list {
		list = append(list, o.old_list[v.Code.Code])
	}

	log.Println("len(old_companys)=", len(o.old_list))
	o.insert.Company = list
	o.insert.Run()
}
