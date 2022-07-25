package company

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/file"
	"github.com/cheolgyu/stockbot/src/fetch/us/us_request"
)

var Addr string = "https://api.nasdaq.com/api/screener/stocks?&download=true&exchange="

type NasdaqCom struct {
	Download bool
}

type NasdaqComItem struct {
	Symbol string
	Name   string
}

const FILE_DIR_US = file.FILE_DIR + "/us" + "/company/"

func (o *NasdaqCom) Request() {

	file.Mkdir([]string{FILE_DIR_US})

	if o.Download {
		for k, _ := range model.Exchanges[model.US] {
			switch k {
			case model.NASDAQ:
				request(Addr+"NASDAQ", model.NASDAQ)
			case model.NYSE:
				request(Addr+"NYSE", model.NYSE)
			case model.AMEX:
				request(Addr+"AMEX", model.AMEX)
			}
		}
	}

}

func request(url string, exchange model.Exchange) {
	fnm := get_fnm(exchange)
	file := file.CreateFile(fnm)

	us_request.HttpNasdaqCom(url, file)
	defer file.Close()
}

func (o *NasdaqCom) GetCompany() []model.Company {
	log.Println("us GetCompany")
	var list []model.Company

	for k, _ := range model.Exchanges[model.US] {
		list = append(list, convert(k)...)
	}

	return list
}

func get_fnm(exchange model.Exchange) string {
	return FILE_DIR_US + model.Exchanges[model.US][exchange].Code
}

func convert(exchange model.Exchange) []model.Company {
	ff := file.File{}

	f := ff.Open(get_fnm(exchange))
	defer f.Close()

	v := make(map[string]map[string][]NasdaqComItem)
	us_request.ConvertNasdaqCom(f, &v)

	arr := v["data"]["rows"]

	var cmp []model.Company
	for _, v := range arr {
		cmp = append(cmp, model.Company{
			Code: model.Code{
				Code: v.Symbol,
				Name: v.Name,
			},
			Country: model.US,
			Market:  model.Exchanges[model.US][exchange].Code,
		})
	}
	return cmp
}
