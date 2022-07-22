package company

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/file"
)

var Addr string = "https://api.nasdaq.com/api/screener/stocks?&download=true&exchange="

type NasdaqCom struct {
}

const FILE_DIR_US = file.FILE_DIR + "/us" + "/company/"

func (o *NasdaqCom) Request() {

	file.Mkdir([]string{FILE_DIR_US})

	fmt.Sprintf("%#v", model.Exchanges[model.US])
	log.Println(model.Exchanges[model.US])

	for k, _ := range model.Exchanges[model.US] {
		fmt.Sprintf("%#v", k)
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

func request(url string, exchange model.Exchange) {
	fnm := FILE_DIR_US + model.Exchanges[model.US][exchange].Code
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
	log.Println("us Save()")
	var list []model.Company

	//파싱하자
	// function  o.나스닥 interface{} 에서 []company 로

	return list
}
