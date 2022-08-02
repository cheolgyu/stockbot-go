package price

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"log"
	"strings"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/file"
)

type NaverChart struct {
	StartDate string
	EndDate   string
	model.Code

	url      string
	fnm      string
	Openings map[int]int
}

const PRICE_DEFAULT_START_DATE = "19560303"
const URL_PRICE = "https://api.finance.naver.com/siseJson.naver?symbol=%s&requestType=1&startTime=%s&endTime=%s&timeframe=day"
const FILE_DIR_PRICE = file.FILE_DIR + "/kr/price/"
const FILE_DIR_MARKET = file.FILE_DIR + "/kr/market/"

func (o *NaverChart) ready() {

	file.Mkdir([]string{FILE_DIR_PRICE, FILE_DIR_MARKET})

	o.fnm = FILE_DIR_PRICE + o.Code.Code

	if o.StartDate == "" {
		o.StartDate = PRICE_DEFAULT_START_DATE
	}

	o.url = fmt.Sprintf(URL_PRICE, o.Code.Code, o.StartDate, o.EndDate)
	o.Openings = make(map[int]int)

}

func (o *NaverChart) GetResult(downlad bool) ([]model.PriceMarket, error) {
	var err error = nil

	o.ready()
	if downlad {
		err_down := o.Download()
		if err_down != nil {
			log.Fatalln(err_down)
			return nil, err_down
		}
	}
	res, err := o.Parse()

	return res, err
}

func (o *NaverChart) Parse() ([]model.PriceMarket, error) {
	var res []model.PriceMarket

	file, err := os.Open(o.fnm)
	if err != nil {
		log.Println("파일열기 에러")
		log.Fatal(err)
		panic(err)
	}

	reader := bufio.NewReader(file)
	idx := 0
	for {
		idx++
		line, isPrefix, err := reader.ReadLine()
		if isPrefix || err != nil {
			break
		}
		str := string(line)
		if idx != 2 {
			var re_str = strings.Replace(str, "[", "", -1)
			re_str = strings.Replace(re_str, "]", "", -1)
			re_str = strings.Replace(re_str, "\"", "", -1)
			re_str = strings.Replace(re_str, " ", "", -1)

			if strings.Contains(re_str, ",") {
				arr := strings.Split(re_str, ",")
				arr[0] = strings.Replace(arr[0], " ", "", -1)
				//dd, e := strconv.ParseInt(arr[0], 0, 64)
				// if e != nil {

				// 	log.Printf("??....%v..", arr[0])
				// 	//panic(e)
				// 	return nil, e
				// }
				// ddd, e := strconv.ParseInt(config.PRICE_DEFAULT_START_DATE, 0, 64)
				// if e != nil {
				// 	return nil, e
				// }

				//if dd > ddd {
				p := stringToPrice(o.Code.Code, re_str)
				o.Openings[p.DateInfo.Dt] = p.DateInfo.Dt
				res = append(res, p)

				//}
			}

		}

	}

	file.Close()

	return res, err

}

func (o *NaverChart) Download() error {
	req, err := http.NewRequest("GET", o.url, nil)
	if err != nil {
		log.Println("Download NewRequest 에러")
		log.Fatal(err)
		return err
	}

	client := &http.Client{
		Timeout: 3 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Download Do 에러")
		log.Fatal(err)
		return err
	}

	out, err := os.Create(o.fnm)
	if err != nil {
		log.Println("Download os.Create 에러")
		log.Fatal(err)
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Download io.Copy 에러")
		log.Fatal(err)
		return err
	}

	out.Close()
	resp.Body.Close()

	return err
}

func parseUint(str string) (int, error) {
	// 08 일경우 오류 발생.
	res, err := strconv.Atoi(str)
	return int(res), err
}

func stringToPrice(code string, str string) model.PriceMarket {
	p := model.PriceMarket{
		Code: code,
	}
	arr := strings.Split(str, ",")
	var s0 = arr[0]

	if res, err := parseUint(s0); err == nil {
		p.DateInfo = model.NewDateInfo(res)
	} else if err != nil {
		panic(err)
	}

	p.OP = model.ParsePrice(arr[1])
	p.HP = model.ParsePrice(arr[2])
	p.LP = model.ParsePrice(arr[3])
	p.CP = model.ParsePrice(arr[4])
	p.Vol = model.ParseVol(arr[5])

	str_fr := strings.Replace(arr[6], ",", "", -1)
	if str_fr == "" {
		str_fr = "0"
	}
	p.FBR = str_fr
	return p
}
