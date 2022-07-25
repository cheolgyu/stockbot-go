package company

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/file"
)

const URL_DETAIL_CODE = "http://data.krx.co.kr/comm/fileDn/GenerateOTP/generate.cmd"
const URL_DETAIL_DATA = "http://data.krx.co.kr/comm/fileDn/download_excel/download.cmd"
const URL_DETAIL_PARAMS = "mktId=ALL&share=1&csvxls_isNo=false&name=fileDown&url=dbms/MDC/STAT/standard/MDCSTAT01901"
const URL_STATE_CODE = "http://data.krx.co.kr/comm/fileDn/GenerateOTP/generate.cmd"
const URL_STATE_DATA = "http://data.krx.co.kr/comm/fileDn/download_excel/download.cmd"
const URL_STATE_PARAMS = "mktId=ALL&share=1&csvxls_isNo=false&name=fileDown&url=dbms/MDC/STAT/standard/MDCSTAT02001"

const FILE_DIR_COMPANY_DETAIL = file.FILE_DIR + "/kr/company_detail/"
const FILE_DIR_COMPANY_STATE = file.FILE_DIR + "/kr/company_state/"
const FILE_NAME_COMPANY_DETAIL = "company_detail.xlsx"
const FILE_NAME_COMPANY_STATE = "company_state.xlsx"
const FILE_COMPANY_DETAIL = FILE_DIR_COMPANY_DETAIL + FILE_NAME_COMPANY_DETAIL
const FILE_COMPANY_STATE = FILE_DIR_COMPANY_STATE + FILE_NAME_COMPANY_STATE

func init() {
	file.Mkdir([]string{FILE_DIR_COMPANY_DETAIL, FILE_DIR_COMPANY_STATE})
}

type Req_krx struct {
	Download bool
	Req_krx_type

	urlCode     string
	urlData     string
	code        string
	codeReqBody string
	saveNm      string
}

type Req_krx_type string

const (
	COMPANY_DETAIL Req_krx_type = "company_detail"
	COMPANY_STATE  Req_krx_type = "company_state"
)

func (o *Req_krx) Request() {
	o.Run()
}

func (o *Req_krx) GetCompany() []model.Company {
	var list []model.Company

	convert := Convert{}
	convert.Run()

	for _, v := range convert.List {
		list = append(list, v)
	}

	return list
}

func (o *Req_krx) init() {

	if o.Req_krx_type == COMPANY_DETAIL {
		o.saveNm = FILE_COMPANY_DETAIL
		o.urlCode = URL_DETAIL_CODE
		o.urlData = URL_DETAIL_DATA
		o.codeReqBody = URL_DETAIL_PARAMS
	} else if o.Req_krx_type == COMPANY_STATE {
		o.saveNm = FILE_COMPANY_STATE
		o.urlCode = URL_STATE_CODE
		o.urlData = URL_STATE_DATA
		o.codeReqBody = URL_STATE_PARAMS
	}
}
func (o *Req_krx) Run() {

	if o.Download {
		detail := Req_krx{Req_krx_type: COMPANY_DETAIL}
		detail.init()
		detail.code = detail.down_code()
		detail.down_file()
		state := Req_krx{Req_krx_type: COMPANY_STATE}
		state.init()
		state.code = state.down_code()
		state.down_file()
	}
}

func (o *Req_krx) down_file() {
	// 파일명
	file := file.CreateFile(o.saveNm)

	reqBody := bytes.NewBufferString("code=" + o.code)
	resp, err := http.Post(o.urlData, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	log.Println("filenm=", o.saveNm, ",size=", size)

	defer file.Close()

	if err != nil {
		panic(err)
	}

}

func (o *Req_krx) down_code() string {
	reqBody := bytes.NewBufferString(o.codeReqBody)
	resp, err := http.Post(o.urlCode, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		str := string(respBody)
		log.Fatalln(str)
	}
	var str_resp = string(respBody)

	return str_resp
}
