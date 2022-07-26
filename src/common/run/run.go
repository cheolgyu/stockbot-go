package run

import (
	"flag"
	"log"
	"strings"

	"github.com/cheolgyu/stockbot/src/common/mlog"
	"github.com/cheolgyu/stockbot/src/common/model"
)

type Run struct {
	Download bool
	model.Country
	mlog.LOG
}

func (o *Run) RunStart() {
	o.Info(o.Country, ",Start")
}
func (o *Run) RunEnd() {
	o.Info(o.Country, ",End")
}

func (o *Run) SetCountry() {

	countryPtr := flag.String("country", "kr", "input country value")
	low_country := strings.ToLower(*countryPtr)
	var country model.Country
	for k, v := range model.Countrys {
		if k == low_country {
			country = v
		}
	}
	o.Country = country

	log.Println("set country : ", string(country))
}
