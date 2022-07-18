package mlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func send(inp LOG) {
	go inp.reqeust()
	go check()
}

type LOG struct {
	Who     string
	What    string
	How     string
	HowWhat string `bson:"how_what"  json:"how_what" `
	Content interface{}
	When    time.Time
	Where   string
}

const LOG_WHAT_INFO = "info"
const LOG_WHAT_ERROR = "error"
const LOG_HOW_WHAT_START = "start"
const LOG_HOW_WHAT_END = "end"

type WHO string

const (
	Fetch     WHO = "fetch"
	AggVol    WHO = "agg_vol"
	LineBound WHO = "line_bound"
	LineYmxb  WHO = "line_ymxb"
	Ticker    WHO = "ticker"
)

var WaitLog sync.WaitGroup
var resp_ch chan *http.Response

func init() {
	resp_ch = make(chan *http.Response)
}
func Info(who WHO, v ...any) {
	l := LOG{
		Who:     string(who),
		What:    LOG_WHAT_INFO,
		When:    time.Now(),
		Content: fmt.Sprint(v...)}

	send(l)
}

func Err(who WHO, v ...any) {
	l := LOG{
		Who:     string(who),
		What:    LOG_WHAT_ERROR,
		When:    time.Now(),
		Content: fmt.Sprint(v...)}
	send(l)
}

func check() {
loop:
	for {
		select {
		case r := <-resp_ch:
			log.Println(r.StatusCode)
			if r.StatusCode == 200 {
				log.Println("log요청.상태코드 정상")
			} else {
				log.Println("log요청.상태코드 !정상=", r.StatusCode)
			}
			break loop
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("^^^^")
		case <-time.After(1 * time.Millisecond):
			fmt.Printf("@")
		}
	}
}

var log_server_url = "http://localhost:5000/logging"

const contentType = "application/json"

func (o *LOG) reqeust() {

	logBytes, err := json.Marshal(o)
	if err != nil {
		log.Fatalln(err)

	}
	resp, err := http.Post(log_server_url, contentType, bytes.NewBuffer(logBytes))
	if err != nil {
		log.Fatalln(err)

	}
	resp_ch <- resp
	println(resp.StatusCode)
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
	}

}
