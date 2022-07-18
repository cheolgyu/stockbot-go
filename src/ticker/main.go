package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/common/model"
)

var marekts []MarketInfo

func init() {
	set_markets()
}
func set_markets() {
	marekts = append(marekts, MarketInfo{
		Country:  model.KR,
		Name:     "한국거래소",
		OpenHour: 0,
		OpenMin:  00,
		ClosHour: 6,
		ClosMin:  00,
	})
	marekts = append(marekts, MarketInfo{
		Country:  model.US,
		Name:     "뉴욕증권거래소, 나스닥",
		OpenHour: 14,
		OpenMin:  30,
		ClosHour: 21,
		ClosMin:  00,
	})
}

func main() {
	if len(os.Args) > 0 {

		pwd := exec.Command("pwd")
		o, err := pwd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(o))

		exec_kr()
	} else {
		m := ment{}
		m.start()
	}
}

//"2006-01-02_15_04_05"
const ment_execd_key_format = "2006-01-02"

type ment struct {
	ticker time.Ticker
	execed map[string]bool
}

func (o *ment) start() {

	ticker := time.NewTicker(time.Second * 60)
	o.ticker = *ticker
	o.execed = map[string]bool{}

	o.loop()
}
func (o *ment) loop() {
	for t := range o.ticker.C {
		utc_t := t.UTC()
		execed_date_key := utc_t.Format(ment_execd_key_format)

		for _, v := range marekts {
			execed_key := fmt.Sprintf("%v", v.Country) + "_" + execed_date_key
			exec_time := time.Date(utc_t.Year(), utc_t.Month(), utc_t.Day(), v.ClosHour, v.ClosMin, 0, 0, time.UTC)

			weekday := fmt.Sprintf("%v", utc_t.Weekday())
			var is_work_day = weekday == "Saturday" || weekday == "Sunday"

			// 시간지남
			case1 := exec_time.After(utc_t)
			// 하루에 한번만
			case2 := o.execed[execed_key]
			// 주말은 쉬자
			case3 := is_work_day

			// p := fmt.Println
			// debug_time := time.Date(utc_t.Year(), utc_t.Month(), utc_t.Day(), utc_t.Hour(), utc_t.Minute()+1, utc_t.Second(), 0, time.UTC)
			// debug_case1 := debug_time.After(utc_t)
			// debug_case2 := !o.execed[execed_key]
			// debug_case3 := !is_work_day
			// p("debug_case1:", debug_case1)
			// p("debug_case2:", debug_case2)
			// p("debug_case3:", debug_case3)

			if case1 && case2 && case3 {
				o.execed[execed_key] = true
				go start(v.Country)
			}
		}
	}
}

//https://ko.wikipedia.org/wiki/증권거래소
type MarketInfo struct {
	model.Country
	Name string
	//UTC기준
	OpenHour int
	//UTC기준
	OpenMin int
	//UTC기준
	ClosHour int
	//UTC기준
	ClosMin int
}

func start(c model.Country) {
	switch c {
	case model.KR:
		exec_kr()
	case model.US:
		exec_us()
	}
}
func exec_kr() {
	mlog.Info(mlog.Ticker, "exec_kr")
	//회사정보
	command("./fetch")
	//가격정보
	command("./asmb_line")
	//가격정보	bound
	//가격정보	bound	y=xm+b
	command("./asmb_agg")
	//가격정보	agg_vol

}
func command(name string) {
	mlog.Info(mlog.Ticker, "command:", name)
	cmd := exec.Command(name)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func exec_us() {
	mlog.Info(mlog.Ticker, "exec_us")
}
