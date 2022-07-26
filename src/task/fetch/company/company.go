package company

import (
	"context"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	kr_company "github.com/cheolgyu/stockbot/src/fetch/country/kr/company"
	us_company "github.com/cheolgyu/stockbot/src/fetch/country/us/company"
	"go.mongodb.org/mongo-driver/mongo"
)

type Crawling interface {
	Request()
	GetCompany() []model.Company
}

var client *mongo.Client

func Run(country model.Country, download bool) {

	client, _ = common.Connect()
	defer client.Disconnect(context.TODO())

	current := SelectMapCompany(client, country)

	var crawling Crawling
	switch country {
	case model.KR:
		crawling = &kr_company.Req_krx{Download: download}
	case model.US:
		crawling = &us_company.NasdaqCom{Download: download}
	}

	crawling.Request()
	incoming := crawling.GetCompany()
	merge_list := merge(country, current, incoming)

	Insert(client, merge_list)
	doc.UpdateNoteOne(doc.GetNoteField(country, doc.PRICE_UPDATE))
}

func merge(country model.Country, current map[string]model.Company, incoming []model.Company) []model.Company {

	var list []model.Company

	for _, v := range incoming {
		current_item, exist := current[v.Code.Code]
		if exist {
			list = append(list, merge_accept_incoming(current_item, v))
		} else {
			list = append(list, merge_accept_new_incoming(country, v))
		}
	}

	return list
}

func merge_accept_incoming(current model.Company, incoming model.Company) model.Company {
	current.Detail = incoming.Detail
	current.State = incoming.State
	return current
}

func merge_accept_new_incoming(country model.Country, incoming model.Company) model.Company {
	incoming.Country = country
	return incoming
}
