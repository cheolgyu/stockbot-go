package price

// func startEndDate() (client *mongo.Client) {

// 	coll := client.Database(doc.DB_PUB).Collection(doc.DB_PUB_COLL_NOTE)
// 	opts := options.Find().SetProjection(bson.M{doc.DB_PUB_COLL_NOTE_PRICE_UPDATED_KR: 1})
// 	cursor, err := coll.Find(ctx, bson.M{}, opts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer cursor.Close(ctx)
// 	var data []bson.M
// 	cursor.All(ctx, &data)
// 	if len(data) == 0 {
// 		log.Println("kr_price_updated_date 신규")
// 		result, err := coll.InsertOne(ctx, bson.M{doc.DB_PUB_COLL_NOTE_PRICE_UPDATED_KR: config.PRICE_DEFAULT_START_DATE})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println("kr_price_updated_date 신규:저장후 id  ", result.InsertedID)
// 		fmt.Println("1.start=", config.PRICE_DEFAULT_START_DATE)
// 	} else {

// 		start = data[0][doc.DB_PUB_COLL_NOTE_PRICE_UPDATED_KR].(string)
// 		log.Println("kr_price_updated_date 존재하는 값 ", start)
// 	}

// 	end = time.Now().Format(config.PRICE_DATE_FORMAT)
// 	fmt.Println(start, end)

// 	return start, end
// }
