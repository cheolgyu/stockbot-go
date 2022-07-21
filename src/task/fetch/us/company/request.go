package company

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cheolgyu/stockbot/src/common/model"
)

var Addr string = "https://api.nasdaq.com/api/screener/stocks?&download=true&exchange="

type NasdaqCom struct {
	NADAQ interface{}
	NYSE  interface{}
	AMEX  interface{}
}

func (o *NasdaqCom) Request() {
	for k, _ := range model.Exchanges[model.US] {
		switch k {
		case model.NASDAQ:
			o.NADAQ = request(Addr + "NASDAQ")
		case model.NYSE:
			o.NYSE = request(Addr + "NYSE")
		case model.AMEX:
			o.AMEX = request(Addr + "AMEX")
		}
	}
}

func request(url string) interface{} {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//필요시 헤더 추가 가능
	req.Header.Add("cache-control", "0")
	req.Header.Add("accept", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var target interface{}
	// 결과 출력
	json.NewDecoder(resp.Body).Decode(target)
	return target
}

func (o *NasdaqCom) GetCompany() []model.Company {
	log.Println("us Save()")
	var list []model.Company

	//파싱하자
	// function  o.나스닥 interface{} 에서 []company 로

	return list
}
