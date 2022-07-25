package price

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/file"
	"github.com/cheolgyu/stockbot/src/fetch/us/us_request"
)

//https://api.nasdaq.com/api/quote/AACG/chart?assetclass=stocks&fromdate=1888-07-25&todate=2022-07-25
const URL string = "https://api.nasdaq.com/api/quote/%s/chart?assetclass=stocks&fromdate=%s&todate=%s"
const FILE_DIR_US = file.FILE_DIR + "/us" + "/price/"
const fromdate string = "1800-01-01"

var todate string

func init() {
	file.Mkdir([]string{FILE_DIR_US})
	todate = time.Now().UTC().Format("2006-01-02")
}

type RequestNasdaqCom struct {
	model.Code
	StartDate string
}
type RespNasdaqCom struct {
	Data struct {
		Chart []struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z struct {
				High     string `json:"high"`
				Low      string `json:"low"`
				Open     string `json:"open"`
				Close    string `json:"close"`
				Volume   string `json:"volume"`
				DateTime string `json:"dateTime"`
				Value    string `json:"value"`
			}
		}
	}
}

func (o *RequestNasdaqCom) GetResult(downlad bool) ([]model.PriceMarket, error) {
	var res []model.PriceMarket

	url := o.getUrl()

	var file *os.File
	defer file.Close()

	if downlad {
		file = o.getSaveFile()
		us_request.HttpNasdaqCom(url, file)
	} else {
		file = o.getFile()
	}

	res = o.convert(file)

	return res, nil
}
func (o *RequestNasdaqCom) convert(f *os.File) []model.PriceMarket {

	var v5 RespNasdaqCom
	us_request.ConvertNasdaqCom(f, &v5)

	var list []model.PriceMarket
	for _, v := range v5.Data.Chart {
		list = append(list, model.PriceMarket{
			Code:     o.Code.Code,
			DateInfo: model.NewDateInfo(cvt_dt(v.Z.DateTime)),
			OP:       cvt_price(v.Z.Open),
			CP:       cvt_price(v.Z.Close),
			LP:       cvt_price(v.Z.Low),
			HP:       cvt_price(v.Z.High),
			Vol:      cvt_vol(v.Z.Volume),
		})
	}
	return list
}
func (o *RequestNasdaqCom) getUrl() string {
	fdt := fromdate
	if o.StartDate != "" {
		fdt = fmt.Sprintf("%s-%s-%s", o.StartDate[:4], o.StartDate[4:6], o.StartDate[6:8])
	}
	return fmt.Sprintf(URL, o.Code.Code, fdt, todate)
}
func (o *RequestNasdaqCom) getSaveFile() *os.File {
	return file.CreateFile(FILE_DIR_US + o.Code.Code)
}
func (o *RequestNasdaqCom) getFile() *os.File {
	f := file.File{}
	return f.Open(FILE_DIR_US + o.Code.Code)
}

// "07/25/2012" to int 20120725
func cvt_dt(dt string) int {
	ys := dt[6:10]
	ms := dt[:2]
	ds := dt[3:5]
	s := ys + ms + ds

	res, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return res
}

// "3.87" to 3.87
func cvt_price(myString string) float32 {
	value, err := strconv.ParseFloat(myString, 32)
	if err != nil {
		panic(err)
	}
	float := float32(value)
	return float
}

//  "2,519" to 2519
func cvt_vol(myString string) int {

	s := strings.Replace(myString, ",", "", -1)
	value, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return value
}
