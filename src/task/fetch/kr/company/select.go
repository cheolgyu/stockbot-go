package company

import (
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
)

func company_map() map[string]model.Company {
	list := doc.GetCompany()
	company_map := make(map[string]model.Company)

	for _, v := range list {
		company_map[v.Code.Code] = v
	}

	return company_map
}
