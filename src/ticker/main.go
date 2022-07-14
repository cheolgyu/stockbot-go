package main

import (
	"fmt"
	"time"

	"github.com/cheolgyu/stockbot/src/common/model"
)

var mlog model.LOG

const WHO = "ticker"

func init() {
	mlog.Who = WHO

}

func main() {

	for i := 0; i < 10; i++ {
		print("!")
	}
	mlog.Log("ㅇㅇㅇㅇㅇㅇㅇㅇㅇㅇㅇㅇㅇ")

	for {
		select {

		case <-time.After(1 * time.Millisecond):
			fmt.Printf(".")
		}
	}
}
