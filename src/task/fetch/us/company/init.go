package company

func init() {
	insert_market()
}

func insert_market() {
	// client, ctx := common.Connect()
	// defer client.Disconnect(ctx)
	// coll := client.Database(doc.DB_PUB).Collection(doc.DB_PUB_COLL_MARKET)

	// for k, v := range model.ExchangeKr {
	// 	if k != model.KONEX {
	// 		filter := bson.M{"code": v.Code, "country": model.KR}
	// 		cnt, err := coll.CountDocuments(context.TODO(), filter)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		if cnt == 0 {
	// 			cmp := model.Market{
	// 				Code: model.Code{
	// 					Code: v.Code,
	// 					Name: v.Name,
	// 				},
	// 				Country: model.KR,
	// 			}
	// 			_, err := coll.InsertOne(context.TODO(), cmp)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 	}
	// }

}
