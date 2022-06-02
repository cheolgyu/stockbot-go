package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// x1,y1 좌표가 bound_point
//x2,y2 는 현재 계산일
type Point struct {
	Code_id    int
	Price_type int
	X1         uint
	Y1         float32
	X2         uint
	Y2         float32
	X_tick     uint
	Y_minus    float32
	Y_Percent  float32
}

func (o *Point) Set(x1 uint, y1 float32, x2 uint, y2 float32, x_tick uint) (err error) {
	// if x_tick > 1 {
	// 	x_tick = x_tick - 1
	// }

	o.X1 = x1
	o.Y1 = y1
	o.X2 = x2
	o.Y2 = y2
	o.X_tick = x_tick
	o.Y_minus = o.Y2 - o.Y1
	o.Y_Percent = float32(float64(o.Y_minus / o.Y2 * 100))

	test := fmt.Sprintf("%v", o.Y_Percent)
	if strings.Contains(test, "Inf") {
		o.Y_Percent = 0
	}
	if x1 > x2 {
		txt := fmt.Sprintf("Point= %v \n", o)
		log.Println(txt)
		err = errors.New("rebound 규칙1. 에 어긋남. x1 < x2")
	}
	return err
}

func StringToPoint(str string) Point {
	item := Point{}

	arr := strings.Split(str, ",")
	x1, _ := strconv.ParseUint(arr[0], 0, 32)
	item.X1 = uint(x1)

	y1, _ := strconv.ParseFloat(arr[1], 32)
	item.Y1 = float32(y1)

	x2, _ := strconv.ParseUint(arr[2], 0, 32)
	item.X2 = uint(x2)

	y2, _ := strconv.ParseFloat(arr[3], 32)
	item.Y2 = float32(y2)

	y_min, _ := strconv.ParseFloat(arr[4], 32)
	item.Y_minus = float32(y_min)

	y_per, _ := strconv.ParseFloat(arr[5], 32)
	item.Y_Percent = float32(y_per)

	x_tick, _ := strconv.ParseUint(arr[6], 0, 32)
	item.X_tick = uint(x_tick)

	return item
}
