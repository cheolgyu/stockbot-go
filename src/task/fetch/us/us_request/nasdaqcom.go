package us_request

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HttpNasdaqCom(url string, f *os.File) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//필요시 헤더 추가 가능
	req.Header.Add("cache-control", "0")
	req.Header.Add("accept", "application/json")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	size, err := io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}
	log.Println("filenm=", f.Name(), ",size=", size)
}

func ConvertNasdaqCom(f *os.File, convertStuct any) {

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, convertStuct)
	// iferr 금지
	// if err != nil {
	// 	panic(err)
	// }

}
