package model

type Bound struct {
	Direction
	Duration uint
	P1       Point
	P2       Point
}
type Direction int

const (
	// nil
	_ Direction = iota
	//감소
	Decrease Direction = iota
	//유지
	Constant Direction = iota
	//증가
	Increase Direction = iota
)

var Directions_arr_int = [...]int{
	0,
	1,
	2,
	3,
}

var Directions_arr_string = [...]string{
	"_",
	"Decrease",
	"Constant",
	"Increase",
}

func (o Direction) String() string { return Directions_arr_string[o-1] }
func (o Direction) Int() int       { return Directions_arr_int[o-1] }
