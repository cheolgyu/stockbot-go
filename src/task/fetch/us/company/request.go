package company

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/file"
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
	defer file.Close()

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
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println("filenm=", fnm, ",size=", size)

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
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var v map[string]map[string][]NasdaqComItem
	json.Unmarshal(data, &v)
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
