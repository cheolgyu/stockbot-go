package main

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/kr/company"
)

func main() {
	log.Println("i am fetch")
	kr_company()
}

func kr_company() {
	old := company.SelectAll()
	log.Println("len(old_companys)=", len(old))
	request_krx := company.Req_krx{}
	request_krx.Run()

	krx_convert := company.Convert{Old: old}
	krx_convert.Run()

	var list []model.Company
	for _, v := range old {
		list = append(list, old[v.Code])
	}

	log.Println("len(old_companys)=", len(old))
	kr_insert := company.Insert{Company: list}
	kr_insert.Run()
}
