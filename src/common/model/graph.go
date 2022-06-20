package model

/*

 */
type Point struct {
	X int
	Y float32
}

type PriceType int

//http://www.gisdeveloper.co.kr/?p=2464
const (
	Open PriceType = 1 + iota
	Close
	Low
	High
)

var PriceTypes = [...]string{
	"Open",
	"Close",
	"Low",
	"High",
}

var PriceTypes_arr = [...]PriceType{
	Open,
	Close,
	Low,
	High,
}

var PriceTypes_DB_Field = [...]string{
	"op",
	"cp",
	"lp",
	"hp",
}

func (o PriceType) String_DB_Field() string { return PriceTypes_DB_Field[o-1] }
func (o PriceType) String() string          { return PriceTypes[o-1] }
