package price

import (
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
)

type Run struct {
	companys []model.Company
	start    string
	end      string
}

/*
1. 종목코드 조회
2. 종목코드로 가격데이터 다운로드
3. 가격데이터 저장 및 company의 마지막가격일자 갱신
*/
func (o *Run) Run() {
	o.companys = doc.GetCompanyCodes()

	start, end := StartEndDate()
	o.start = start
	o.end = end

}
