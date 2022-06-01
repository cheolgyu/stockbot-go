package price

import (
	"log"

	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
)

type Run struct {
	code     []model.Code
	start    string
	end      string
	download naverChart
	insert   Insert
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 pub.note 마지막가격일자 갱신
*/
func (o *Run) Run() {

	o.code = doc.GetCodes()
	o.start, o.end = startEndDate()

	for _, v := range o.code {
		o.download = naverChart{
			startDate: o.start,
			endDate:   o.end,
			Code:      v,
		}
		list, err := o.download.Run()
		if err != nil {
			log.Panic(err)
		}
		list
	}

}
