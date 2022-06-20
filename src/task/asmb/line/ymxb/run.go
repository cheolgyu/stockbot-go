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

/*
	y=mx+b,
	p1은 저고종시가의 마지막 바운스점,
	p2는 저고종시가의 마지막 점.

	p1~3는 실제 가격이고 tp1~3은 추상화한 가격,
	tp의 x는 일자를 1칸으로 추상화하였고,
	tp의 y는 호가를 1칸으로 추상화하였음.
*/
type ymxb struct {
	code        string
	market_code model.Market
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

func Run() {
	client, ctx := common.Connect()
	defer client.Disconnect(ctx)
	collection_boud_point := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_BOUND_POINT)
	collection_price := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_PRICE)
	collection_ymxb_quote_unit := client.Database(doc.DB_DATA).Collection(doc.DB_DATA_COLL_YMXB_QUOTE_UNIT)

	company := doc.GetCompany()

	for _, c := range company {
		market_code, err := model.String2Market(c.Market)
		if err != nil {
			log.Panic(err.Error())
		}

		for _, v := range model.PriceTypes_arr {
			i := ymxb{
				code:        c.Code.Code,
				market_code: market_code,
				price_type:  v,
			}

			i.setP2_last_price(collection_price)
			i.setP1_last_bound(collection_boud_point)

			// p1점이 없는 경우
			if i.p1.X != 0 {
				i.setM(collection_price, collection_ymxb_quote_unit)
				panic("저장하기 만들기")
			}

		}

	}
}

func errHandler(inp error, v ...interface{}) {
	if inp != nil {
		fmt.Printf("%#v \n", v...)
		log.Fatalln(inp)
	}

}

func (o *ymxb) setP1_last_bound(coll *mongo.Collection) {

	projection := bson.M{"_id": 0, "x": "$p2.x", "y": "$p2.y"}
	filter := bson.M{"code": o.code, "dt": bson.M{"$lt": o.p2.X}, "price_type": o.price_type}
	opts := options.FindOne().SetSort(bson.M{"p2.x": -1}).SetProjection(projection)

	err := coll.FindOne(context.TODO(), filter, opts).Decode(&o.p1)
	if err != nil && mongo.ErrNoDocuments != err {
		log.Println("????", err)
		errHandler(err, o, "setP1")
	}

}
func (o *ymxb) setP2_last_price(coll *mongo.Collection) {

	projection := bson.M{"_id": 0, "x": "$dt", "y": "$" + o.price_type.String_DB_Field()}
	filter := bson.M{"code": o.code}
	opts := options.FindOne().SetSort(bson.M{"dt": -1}).SetProjection(projection)

	err := coll.FindOne(context.TODO(), filter, opts).Decode(&o.p2)
	errHandler(err, o, "setP2")

}

func (o *ymxb) setM(coll_price *mongo.Collection, coll_ymxb_quote_unit *mongo.Collection) {

	filter := bson.M{"code": o.code, "dt": bson.M{"$lte": o.p1.X}}
	p1x, err := coll_price.CountDocuments(context.TODO(), filter)
	errHandler(err, o, "setM", "o.p1.x")

	filter = bson.M{"code": o.code, "dt": bson.M{"$lte": o.p2.X}}
	p2x, err := coll_price.CountDocuments(context.TODO(), filter)
	errHandler(err, o, "setM", "o.p2.x")

	tp1y := model.Unit_quote{}
	tp2y := model.Unit_quote{}
	filter2 := bson.M{"market": o.market_code, "price": o.p1.Y}
	err = coll_ymxb_quote_unit.FindOne(context.TODO(), filter2).Decode(&tp1y)
	errHandler(err, o, "setM", "o.p1.y")

	filter2 = bson.M{"market": o.market_code, "price": o.p2.Y}
	err = coll_ymxb_quote_unit.FindOne(context.TODO(), filter2).Decode(&tp2y)
	errHandler(err, o, "setM", "o.p2.y")

	fmt.Println(o.code, o.price_type)
	o.tp1x = float64(p1x)
	o.tp2x = float64(p2x)
	o.tp1y = float64(tp1y.Tick)
	o.tp2y = float64(tp2y.Tick)
	log.Println("tp1=", o.tp1x, o.tp1y, ",tp2=", o.tp2x, o.tp2y)

	o.m = model.Round2((o.tp2y - o.tp1y) / (o.tp2x - o.tp1x))
	// y=mx+b
	// y -mx = b
	o.b = float64(o.tp1y + o.m*(-1)*o.tp1x)
	log.Println(" y=  ", o.m, "x +", o.b)
	o.tp3x = o.tp2x + 1
	o.tp3y = o.m*o.tp3x + o.b

	filter2 = bson.M{"market": o.market_code, "tick": o.tp3y}
	opts := options.FindOne()

	tp3 := model.Unit_quote{}
	err = coll_ymxb_quote_unit.FindOne(context.TODO(), filter2, opts).Decode(&tp3)
	errHandler(err, o, "setM", "tp3")

	o.p3.Y = float32(tp3.Price)
	log.Println("p3.y=", o.p3)
}
