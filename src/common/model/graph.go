package model

/*

 */
type Point struct {
	X uint
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

func (o PriceType) String() string { return PriceTypes[(o - 1)] }
