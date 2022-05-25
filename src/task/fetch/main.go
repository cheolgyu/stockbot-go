package main

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/kr"
)

func main() {
	log.Println("i am fetch")
	fetch_kr()
}

func fetch_kr() {
	old := kr.SelectAll()
	log.Println("len(old_companys)=", len(old))
	request_krx := kr.Req_krx{}
	request_krx.Run()

	krx_convert := kr.Convert{Old: old}
	krx_convert.Run()

	var list []model.Company
	for _, v := range old {
		list = append(list, old[v.Code])
	}

	log.Println("len(old_companys)=", len(old))
	kr_insert := kr.Insert{Company: list}
	kr_insert.Run()
}
