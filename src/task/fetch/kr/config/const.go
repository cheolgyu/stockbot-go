package config

import (
	"os"
)

const DOWNLOAD_URL_COMPANY_DETAIL_CODE = "http://data.krx.co.kr/comm/fileDn/GenerateOTP/generate.cmd"
const DOWNLOAD_URL_COMPANY_DETAIL_DATA = "http://data.krx.co.kr/comm/fileDn/download_excel/download.cmd"
const DOWNLOAD_URL_COMPANY_DETAIL_PARAMS = "mktId=ALL&share=1&csvxls_isNo=false&name=fileDown&url=dbms/MDC/STAT/standard/MDCSTAT01901"
const DOWNLOAD_URL_COMPANY_STATE_CODE = "http://data.krx.co.kr/comm/fileDn/GenerateOTP/generate.cmd"
const DOWNLOAD_URL_COMPANY_STATE_DATA = "http://data.krx.co.kr/comm/fileDn/download_excel/download.cmd"
const DOWNLOAD_URL_COMPANY_STATE_PARAMS = "mktId=ALL&share=1&csvxls_isNo=false&name=fileDown&url=dbms/MDC/STAT/standard/MDCSTAT02001"
const DOWNLOAD_URL_PRICE = "https://api.finance.naver.com/siseJson.naver?symbol=%s&requestType=1&startTime=%s&endTime=%s&timeframe=day"
const DOWNLOAD_URL_FIND_MARKET = "https://finance.naver.com/item/coinfo.nhn?code="

const DOWNLOAD_DIR_COMPANY_DETAIL = "data/download/company_detail/"
const DOWNLOAD_DIR_COMPANY_STATE = "data/download/company_state/"
const DOWNLOAD_FILENAME_COMPANY_DETAIL = "company_detail.xlsx"
const DOWNLOAD_FILENAME_COMPANY_STATE = "company_state.xlsx"
const DOWNLOAD_COMPANY_DETAIL = DOWNLOAD_DIR_COMPANY_DETAIL + DOWNLOAD_FILENAME_COMPANY_DETAIL
const DOWNLOAD_COMPANY_STATE = DOWNLOAD_DIR_COMPANY_STATE + DOWNLOAD_FILENAME_COMPANY_STATE

const DOWNLOAD_DIR_PRICE = "data/download/price/"
const DOWNLOAD_DIR_MARKET = "data/download/market/"

const COMPANY_DETAIL = "company_detail"
const COMPANY_STATE = "company_state"

var DownloadCompany bool
var DownloadPrice bool

const FILE_FLAG_APPEND = os.O_RDWR | os.O_CREATE | os.O_APPEND
const FILE_FLAG_TRUNC = os.O_RDWR | os.O_CREATE | os.O_TRUNC

const INFO_NAME_UPDATED = "updated_stock_write"

//custom := now.Format("2006-01-02 15:04:05")
const PRICE_DATE_FORMAT = "20060102"

const PRICE_DEFAULT_START_DATE = "19560303"

const XLSX_SPLIT = "!,_"

var DB_MAX_CONN = 30

var Config map[string]int

func init() {
	DownloadCompany = false
	DownloadPrice = false

}
