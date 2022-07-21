package company

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/kr/config"
)

func (o *Req_krx) Request() {
	req_krx := Req_krx{}
	req_krx.Run()
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

type Req_krx struct {
	Object string

	urlCode     string
	urlData     string
	code        string
	codeReqBody string
	saveNm      string
}

func (o *Req_krx) init() {
	if o.Object == config.COMPANY_DETAIL {
		o.saveNm = config.DOWNLOAD_DIR_COMPANY_DETAIL + config.DOWNLOAD_FILENAME_COMPANY_DETAIL
		o.urlCode = config.DOWNLOAD_URL_COMPANY_DETAIL_CODE
		o.urlData = config.DOWNLOAD_URL_COMPANY_DETAIL_DATA
		o.codeReqBody = config.DOWNLOAD_URL_COMPANY_DETAIL_PARAMS
	} else if o.Object == config.COMPANY_STATE {
		o.saveNm = config.DOWNLOAD_DIR_COMPANY_STATE + config.DOWNLOAD_FILENAME_COMPANY_STATE
		o.urlCode = config.DOWNLOAD_URL_COMPANY_STATE_CODE
		o.urlData = config.DOWNLOAD_URL_COMPANY_STATE_DATA
		o.codeReqBody = config.DOWNLOAD_URL_COMPANY_STATE_PARAMS
	}
}
func (o *Req_krx) Run() {

	if config.DownloadCompany {
		detail := Req_krx{Object: config.COMPANY_DETAIL}
		detail.init()
		detail.code = detail.down_code()
		detail.down_file()
		state := Req_krx{Object: config.COMPANY_STATE}
		state.init()
		state.code = state.down_code()
		state.down_file()
	}
}

func (o *Req_krx) down_file() {
	// 파일명
	file := createFile(o.saveNm)

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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func createFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	checkError(err)
	return file
}
