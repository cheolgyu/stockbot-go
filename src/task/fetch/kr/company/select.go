package company

import (
	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func init() {
	client, _ = common.Connect()
}

func company_map() map[string]model.Company {
	list := doc.GetCompany(client)
	company_map := make(map[string]model.Company)

	for _, v := range list {
		company_map[v.Code.Code] = v
	}

	return company_map
}
