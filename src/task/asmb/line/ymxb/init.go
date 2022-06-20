package ymxb

import (
	"context"
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	create_index()
	if !check() {
		create_index()
		initial_db_unit_coll_ymxb_qout()
	}
	println("ymxb runed")
}

/*
호가테이블 생성하기
*/
func initial_db_unit_coll_ymxb_qout() {
	list := kr_kospi()
	list = append(list, kr_kosdaq_konex(model.KOSDAQ)...)
	list = append(list, kr_kosdaq_konex(model.KONEX)...)
	insert(list)
}

func kr_kospi() []model.Unit_quote {
	var list []model.Unit_quote
	count := 5999999
	tick := 1

	for i := 1; i < count; i++ {
		save := false
		if i < 1000 {
			save = true
		} else if 1000 <= i && i < 5000 && i%5 == 0 {
			save = true
		} else if 5000 <= i && i < 10000 && i%10 == 0 {
			save = true
		} else if 10000 <= i && i < 50000 && i%50 == 0 {
			save = true
		} else if 50000 <= i && i < 100000 && i%100 == 0 {
			save = true
		} else if 100000 <= i && i < 500000 && i%500 == 0 {
			save = true
		} else if 500000 <= i && i%1000 == 0 {
			save = true
		}

		if save {
			list = append(list, model.Unit_quote{
				Market: model.KOSPI,
				Tick:   tick,
				Price:  i,
			})
			tick++
		}
	}
	return list
}
func kr_kosdaq_konex(market model.Market) []model.Unit_quote {
	var list []model.Unit_quote
	count := 5999999
	tick := 1

	for i := 1; i < count; i++ {
		save := false
		if i < 1000 {
			save = true
		} else if 1000 <= i && i < 5000 && i%5 == 0 {
			save = true
		} else if 5000 <= i && i < 10000 && i%10 == 0 {
			save = true
		} else if 10000 <= i && i < 50000 && i%50 == 0 {
			save = true
		} else if 50000 <= i && i < 100000 && i%100 == 0 {
			save = true
		} else if 100000 <= i && i < 500000 && i%100 == 0 {
			save = true
		} else if 500000 <= i && i%100 == 0 {
			save = true
		}

		if save {
			list = append(list, model.Unit_quote{
				Market: market,
				Tick:   tick,
				Price:  i,
			})
			tick++
		}
	}
	return list
}
func insert(list []model.Unit_quote) {
	var newValue []interface{}

	for _, v := range list {
		var i interface{}
		i = v
		newValue = append(newValue, i)
	}

	println("ymxb insert")
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	collection := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_YMXB_QUOTE_UNIT)
	result, err := collection.InsertMany(ctx, newValue)

	if err != nil {
		log.Fatalln("err.Error():", err.Error())
	}
	fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
}
func check() bool {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	count, err := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_YMXB_QUOTE_UNIT).EstimatedDocumentCount(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Estimated number of documents in the ratings collection: %d\n", count)
	return count > 0
}

//https://christiangiacomi.com/posts/mongodb-index-using-go/
func create_index() {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)

	dataPriceCollection := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_YMXB_QUOTE_UNIT)

	/*
			// IndexModel represents a new index to be created.
		type IndexModel struct {
			// A document describing which keys should be used for the index. It cannot be nil. This must be an order-preserving
			// type such as bson.D. Map types such as bson.M are not valid. See https://docs.mongodb.com/manual/indexes/#indexes
			// for examples of valid documents.
			Keys interface{}

			// The options to use to create the index.
			Options *options.IndexOptions
		}
	*/
	mod := mongo.IndexModel{
		Keys:    bson.D{{"market", 1}, {"price", 1}, {"tick", 1}}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := dataPriceCollection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		panic(err.Error())
	}
}
