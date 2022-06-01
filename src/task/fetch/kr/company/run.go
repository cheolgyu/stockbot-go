package company

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/model"
)

type Run struct{}

func (o *Run) Run() {
	old := company_map()
	log.Println("len(old_companys)=", len(old))
	request_krx := Req_krx{}
	request_krx.Run()

	krx_convert := Convert{Old: old}
	krx_convert.Run()

	var list []model.Company
	for _, v := range old {
		list = append(list, old[v.Code.Code])
	}

	log.Println("len(old_companys)=", len(old))
	kr_insert := Insert{Company: list}
	kr_insert.Run()
}
