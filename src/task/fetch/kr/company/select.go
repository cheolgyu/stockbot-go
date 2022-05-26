package company

import (
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/model"
	"github.com/cheolgyu/stockbot/src/fetch/kr"
	"go.mongodb.org/mongo-driver/bson"
)

func SelectAll() map[string]model.Company {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	cursor, err := client.Database(kr.DB_PUB).Collection(kr.DB_PUB_COLL_COMPANY).Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var company_map map[string]model.Company
	company_map = make(map[string]model.Company)

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result model.Company
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		company_map[result.Code] = result
	}

	fmt.Println("select company_map len=", len(company_map))
	return company_map
}
