package ymxb

import (
	"context"
	"fmt"
	"log"

	"github.com/cheolgyu/stockbot/src/common"
	"github.com/cheolgyu/stockbot/src/common/base"
	"github.com/cheolgyu/stockbot/src/common/doc"
	"github.com/cheolgyu/stockbot/src/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection_boud_point *mongo.Collection
var collection_price *mongo.Collection
var collection_ymxb_quote_unit *mongo.Collection
var collection_ymxb *mongo.Collection

type LineYmxb struct {
	base.Run
}

/*
	y=mx+b,
	p1은 저고종시가의 마지막 바운스점,
	p2는 저고종시가의 마지막 점.

	p1~3는 실제 가격이고 tp1~3은 추상화한 가격,
	tp의 x는 일자를 1칸으로 추상화하였고,
	tp의 y는 호가를 1칸으로 추상화하였음.

	+ 현재는 p1은 마지막 바운스고 p2는 마지막점을 type1로 두고
	+ 나중에 type2는 마지막 저점 바운스 2개고
	+ type3는 마지막 고점 바운스 2개 로 만들기.
*/
type ymxb struct {
	code        string
	market_code model.Exchange
	price_type  model.PriceType

	p1 model.Point
	p2 model.Point
	p3 model.Point

	tp1x float64
	tp1y float64
	tp2x float64
	tp2y float64

	tp3x float64
	tp3y float64

	m float64
	b float64
}

func (o *LineYmxb) EXE() {

	client, _ = common.Connect()
	defer client.Disconnect(context.TODO())

	collection_boud_point = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_BOUND_POINT)
	collection_price = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
	collection_ymxb_quote_unit = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_YMXB_QUOTE_UNIT)
	collection_ymxb = client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_YMXB)

	company := doc.GetCompany(client, o.Country)
	var list []interface{}

	for _, c := range company {
		market_code, err := model.ConvertExchanges(o.Country, c.Market)
		errHandler(err, "model.String2Market(c.Market),c.Market= ", c.Market)

		for _, v := range model.PriceTypes_arr {
			i := ymxb{
				code:        c.Code.Code,
				market_code: market_code,
				price_type:  v,
			}

			i.setP2_last_price()
			i.setP1_last_bound()

			// p1점이 없는 경우
			if i.p1.X != 0 {
				i.setM()
				list = append(list, model.Ymxb{
					Code:      i.code,
					YmxbType:  model.YmxbType1,
					PriceType: i.price_type,
					P1:        i.p1,
					P2:        i.p2,
					P3:        i.p3,
					M:         i.m,
					B:         i.b,
				})
			}
		}
	}
	err := collection_ymxb.Drop(context.TODO())
	errHandler(err, "collection_ymxb.Drop")

	mod := mongo.IndexModel{
		Keys:    bson.D{{"code", 1}, {"price_type", 1}, {"ymxb_type", model.YmxbType1}}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}
	_, err = collection_ymxb.Indexes().CreateOne(context.TODO(), mod)
	errHandler(err, "collection_ymxb.Indexes().CreateOne")
	_, err = collection_ymxb.InsertMany(context.TODO(), list)
	errHandler(err, "collection_ymxb.InsertMany")
}

func errHandler(inp error, v ...interface{}) {
	if inp != nil {
		fmt.Printf("%#v \n", v...)
		log.Fatalln(inp)
	}
}

func (o *ymxb) setP1_last_bound() {

	projection := bson.M{"_id": 0, "x": "$p2.x", "y": "$p2.y"}
	filter := bson.M{"code": o.code, "p2.x": bson.M{"$lt": o.p2.X}, "price_type": o.price_type}
	opts := options.FindOne().SetSort(bson.M{"p2.x": -1}).SetProjection(projection)

	err := collection_boud_point.FindOne(context.TODO(), filter, opts).Decode(&o.p1)
	if err != nil && mongo.ErrNoDocuments != err {
		log.Println("????", err)
		errHandler(err, o, "setP1")
	}

}
func (o *ymxb) setP2_last_price() {

	projection := bson.M{"_id": 0, "x": "$dt", "y": "$" + o.price_type.String_DB_Field()}
	filter := bson.M{"code": o.code}
	opts := options.FindOne().SetSort(bson.M{"dt": -1}).SetProjection(projection)

	err := collection_price.FindOne(context.TODO(), filter, opts).Decode(&o.p2)
	errHandler(err, o, "setP2")

}

func (o *ymxb) setM() {
	filter := bson.M{"code": o.code, "dt": bson.M{"$lte": o.p1.X}}
	p1x, err := collection_price.CountDocuments(context.TODO(), filter)
	errHandler(err, o, "setM", "o.p1.x")

	filter = bson.M{"code": o.code, "dt": bson.M{"$lte": o.p2.X}}
	p2x, err := collection_price.CountDocuments(context.TODO(), filter)
	errHandler(err, o, "setM", "o.p2.x")

	tp1y := model.Unit_quote{}
	tp2y := model.Unit_quote{}
	filter2 := bson.M{"market": o.market_code, "price": o.p1.Y}
	err = collection_ymxb_quote_unit.FindOne(context.TODO(), filter2).Decode(&tp1y)
	errHandler(err, o, "setM", "o.p1.y")

	filter2 = bson.M{"market": o.market_code, "price": o.p2.Y}
	err = collection_ymxb_quote_unit.FindOne(context.TODO(), filter2).Decode(&tp2y)
	errHandler(err, o, "setM", "o.p2.y")

	o.tp1x = float64(p1x)
	o.tp2x = float64(p2x)
	o.tp1y = float64(tp1y.Tick)
	o.tp2y = float64(tp2y.Tick)
	//log.Println("tp1=", o.tp1x, o.tp1y, ",tp2=", o.tp2x, o.tp2y)

	o.m = model.Round2((o.tp2y - o.tp1y) / (o.tp2x - o.tp1x))
	// y=mx+b
	// y -mx = b
	o.b = float64(o.tp1y + o.m*(-1)*o.tp1x)
	//log.Println(" y=  ", o.m, "x +", o.b)
	o.tp3x = o.tp2x + 1
	o.tp3y = o.m*o.tp3x + o.b

	filter2 = bson.M{"market": o.market_code, "tick": model.Round0(o.tp3y)}
	opts := options.FindOne()

	tp3 := model.Unit_quote{}
	err = collection_ymxb_quote_unit.FindOne(context.TODO(), filter2, opts).Decode(&tp3)
	errHandler(err, o, "setM", "tp3")

	o.p3.Y = float32(tp3.Price)
	//log.Println("p3.y=", o.p3)
}
