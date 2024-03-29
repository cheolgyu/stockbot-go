package price

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/country/us/us_request"
	"github.com/cheolgyu/stockbot/src/fetch/file"
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
		file.Close()
	}

	file = o.getFile()
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
			OP:       model.ParsePrice(v.Z.Open),
			CP:       model.ParsePrice(v.Z.Close),
			LP:       model.ParsePrice(v.Z.Low),
			HP:       model.ParsePrice(v.Z.High),
			Vol:      model.ParseVol(v.Z.Volume),
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
