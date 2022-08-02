package company

import (
	"fmt"
	"log"
	"strings"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/tealeg/xlsx"
)

const XLSX_SPLIT = "!,_"

type Convert struct {
	List map[string]model.Company
}

func (o *Convert) Run() {
	o.List = make(map[string]model.Company)

	o.run_detail()
	o.run_state()

}

func (o *Convert) run_state() {

	xlFile, err := xlsx.OpenFile(FILE_COMPANY_STATE)
	if err != nil {
		log.Fatalln("run_state 오류발생", err)
	}

	sheet := xlFile.Sheets[0]

	for i := 1; i < sheet.MaxRow; i++ {
		row := sheet.Row(i)
		_, content := rowGet(row)
		state := stringToCompanyState(content)
		tmp := o.List[state.Code.Code]
		tmp.State = state.State
		o.List[state.Code.Code] = tmp
	}
}

func (o *Convert) run_detail() {

	xlFile, err := xlsx.OpenFile(FILE_COMPANY_DETAIL)
	if err != nil {
		log.Fatalln("run_detail 오류발생", err)
	}

	sheet := xlFile.Sheets[0]

	for i := 1; i < sheet.MaxRow; i++ {
		row := sheet.Row(i)
		_, content := rowGet(row)
		detail := stringToCompanyDetail(content)
		o.List[detail.Code.Code] = detail
	}
}

func rowGet(row *xlsx.Row) (string, string) {
	txt_replace := strings.NewReplacer("'", " ")

	str := fmt.Sprintf("%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s"+XLSX_SPLIT+"%s",
		txt_replace.Replace(row.Cells[0].String()),
		txt_replace.Replace(row.Cells[1].String()),
		txt_replace.Replace(row.Cells[2].String()),
		txt_replace.Replace(row.Cells[3].String()),
		txt_replace.Replace(row.Cells[4].String()),
		txt_replace.Replace(row.Cells[5].String()),
		txt_replace.Replace(row.Cells[6].String()),
		txt_replace.Replace(row.Cells[7].String()),
		txt_replace.Replace(row.Cells[8].String()),
		txt_replace.Replace(row.Cells[9].String()),
		txt_replace.Replace(row.Cells[10].String()),
		txt_replace.Replace(row.Cells[11].String()),
	)

	return row.Cells[1].String(), str
}

func stringToCompanyDetail(str string) model.Company {
	arr := strings.Split(str, XLSX_SPLIT)

	cmp := model.Company{}

	cmp_detail := model.CompanyDetail{}
	cmp_detail.Full_code = arr[0]
	cmp.Code.Code = arr[1]
	cmp_detail.Full_name_kr = arr[2]
	cmp.Code.Name = arr[3]

	cmp_detail.Full_name_eng = arr[4]

	cmp_detail.Listing_date = arr[5]
	cmp.Market = strings.ToLower(arr[6])
	cmp_detail.Market = cmp.Market
	cmp_detail.Securities_classification = arr[7]
	cmp_detail.Department = arr[8]
	cmp_detail.Stock_type = arr[9]

	cmp_detail.Face_value = arr[10]
	cmp_detail.Listed_stocks = arr[11]

	cmp.Detail = cmp_detail

	return cmp
}

func stringToCompanyState(str string) model.Company {

	ic := model.Company{}
	o := model.CompanyState{}
	arr := strings.Split(str, XLSX_SPLIT)

	txt_replace := strings.NewReplacer("'", " ")

	ic.Code.Code = txt_replace.Replace(arr[0])
	ic.Code.Name = txt_replace.Replace(arr[1])
	o.Stop = ox(arr[2])
	o.Clear = ox(arr[3])
	o.Managed = ox(arr[4])

	o.Ventilation = ox(arr[5])
	o.Unfaithful = ox(arr[6])
	o.Lack_listed = ox(arr[7])
	o.Overheated = ox(arr[8])

	o.Caution = ox(arr[9])
	o.Warning = ox(arr[10])
	o.Risk = ox(arr[11])

	ic.State = o
	return ic
}

func ox(ox string) bool {

	if ox == "X" {
		return false
	} else {
		return true
	}

}
