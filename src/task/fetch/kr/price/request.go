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
	"github.com/cheolgyu/stockbot/src/fetch/kr/config"
)

type naverChart struct {
	startDate string
	endDate   string
	model.Code

	url      string
	fnm      string
	Openings map[int]int
}

func (o *naverChart) ready() {
	// if o.Code.Code.Code_type == config.Config["stock"] {
	// 	o.fnm = config.DOWNLOAD_DIR_PRICE + o.Code.Code
	// } else if o.Code.Code.Code_type == config.Config["market"] {
	// 	o.fnm = config.DOWNLOAD_DIR_MARKET + o.Code.Code
	// }
	o.fnm = config.DOWNLOAD_DIR_PRICE + o.Code.Code

	if o.startDate == "" {
		o.startDate = config.PRICE_DEFAULT_START_DATE
	}

	o.url = fmt.Sprintf(config.DOWNLOAD_URL_PRICE, o.Code.Code, o.startDate, o.endDate)
	o.Openings = make(map[int]int)

}

func (o *naverChart) Run() ([]model.PriceMarket, error) {
	var err error = nil

	o.ready()
	if config.DownloadPrice {
		err_down := o.Download()
		if err_down != nil {
			log.Fatalln(err_down)
			return nil, err_down
		}
	}
	res, err := o.Parse()

	return res, err
}

func (o *naverChart) Parse() ([]model.PriceMarket, error) {
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
				o.Openings[p.Dt] = p.Dt
				res = append(res, p)

				//}
			}

		}

	}

	file.Close()

	return res, err

}

func (o *naverChart) Download() error {
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
		p.Dt = res
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[1], 32); err == nil {
		p.OP = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[2], 32); err == nil {
		p.HP = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[3], 32); err == nil {
		p.LP = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := strconv.ParseFloat(arr[4], 32); err == nil {
		p.CP = float32(res)
	} else if err != nil {
		panic(err)
	}

	if res, err := parseUint(arr[5]); err == nil {
		p.Vol = res
	} else if err != nil {
		panic(err)
	}

	str_fr := strings.Replace(arr[6], ",", "", -1)
	if str_fr == "" {
		str_fr = "0"
	}
	p.FBR = str_fr
	return p
}
