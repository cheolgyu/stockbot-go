package model

type Bound struct {
	Direction
	P1 Point
	P2 Point
}
type Direction int

const (
	//감소
	Decrease Direction = -1
	//유지
	Constant Direction = 0
	//증가
	Increase Direction = 1
)

var Directions_arr_int = [...]int{
	-1,
	0,
	1,
}

var Directions_arr_string = [...]string{
	"Decrease",
	"Constant",
	"Increase",
}

func (o Direction) String() string { return Directions_arr_string[o-1] }
func (o Direction) Int() int       { return Directions_arr_int[o-1] }
