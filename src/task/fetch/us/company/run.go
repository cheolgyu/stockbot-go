package company

import (
	"context"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/dao"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func Run() {
	client, _ = common.Connect()
	defer client.Disconnect(context.TODO())
	for _, country := range model.Countrys {
		current := dao.SelectMapCompany(client, country)
	}

	// incoming :=
	// merge :=
	// insert

}

func merge(country model.Country, current map[string]model.Company, incoming []model.Company) []model.Company {

	var list []model.Company

	for _, v := range incoming {
		current_item, exist := current[v.Code.Code]
		if exist {
			list = append(list, merge_accept_incoming(current_item, v))
		} else {
			list = append(list, merge_accept_new_incoming(country, current_item))
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
