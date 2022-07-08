package model

type PeriodType int

const (
	//주
	Weeks PeriodType = iota
	//월
	Month
	//분기
	Quarter
)

func (o PeriodType) toField() string {
	str := ""
	if o == Weeks {
		str = "weeks"
	} else if o == Month {
		str = "month"
	} else if o == Quarter {
		str = "quarter"
	}
	return str
}

type AggVolSum struct {
	Code       string
	Year       int
	SumWeeks   map[int]int
	SumMonth   map[int]int
	SumQuarter map[int]int
	Result     map[string]AggVolSumResult
}

func (o *AggVolSum) GetValueOfObservationValueType(ov ObservationValueType) int {
	res := -1

	if ov == WeekMinValue {
		res = o.Result[Weeks.toField()].Min
	} else if ov == WeekMaxValue {
		res = o.Result[Weeks.toField()].Max
	} else if ov == MonthMinValue {
		res = o.Result[Month.toField()].Min
	} else if ov == MonthMaxValue {
		res = o.Result[Month.toField()].Max
	} else if ov == QuarterMinValue {
		res = o.Result[Quarter.toField()].Min
	} else if ov == QuarterMaxValue {
		res = o.Result[Quarter.toField()].Max
	} else {
		panic("ObservationValueType 변환오류")
	}

	return res
}

type AggVolSumResult struct {
	Min       int
	Max       int
	Avg       int
	Avg_under []int
	Avg_upper []int
}

// 3. 코드별 연도별 가격데이터에서 연도의 최소거래량인 주,월,분기 최대거래량인 주,월,분기, 평균, 평균이하이상을 구한다.
func (o *AggVolSum) Calculate() {
	o.Result = make(map[string]AggVolSumResult)

	o.Result[Weeks.toField()] = calculateByPeriodType(o.SumWeeks)
	o.Result[Month.toField()] = calculateByPeriodType(o.SumMonth)
	o.Result[Quarter.toField()] = calculateByPeriodType(o.SumQuarter)

	//fmt.Println("%v", o.Result)
}

func calculateByPeriodType(items map[int]int) (res AggVolSumResult) {
	min_idx, max_idx := 1, -1
	min_tmp, max_tmp := -1, -1
	sum_vol := 0
	cnts := 0

	for k, v := range items {
		sum_vol += v
		cnts++

		if max_tmp < v {
			max_tmp = v
			max_idx = k
		}
		if min_tmp > v || min_tmp == -1 {
			min_tmp = v
			min_idx = k
		}

	}

	avg_vol := sum_vol / cnts
	var under_avg []int
	var upper_avg []int

	for k, v := range items {
		if v > avg_vol {
			upper_avg = append(upper_avg, k)
		} else {
			under_avg = append(under_avg, k)
		}
	}

	res.Min = min_idx
	res.Max = max_idx
	res.Avg = avg_vol
	res.Avg_under = under_avg
	res.Avg_upper = upper_avg

	//log.Println("%v", res)

	return res
}

type AggVol struct {
	Code   string
	Result map[string]AggVolStatisticBasic
}

type AggVolStatisticBasic struct {

	// 데이터수
	DataCnt int
	// 데이터 key:년도 value:관측값
	Data map[int]int
	// 평균
	Average int
	// 분산
	Variance float64
	// 표준편차
	StandardDeviation float64
}

// 관찰값종류
type ObservationValueType int

func (o ObservationValueType) ToField() string {
	str := ""
	if o == WeekMinValue {
		str = "weeks_min"
	} else if o == WeekMaxValue {
		str = "weeks_max"
	} else if o == MonthMinValue {
		str = "month_min"
	} else if o == MonthMaxValue {
		str = "month_max"
	} else if o == QuarterMinValue {
		str = "quarter_min"
	} else if o == QuarterMaxValue {
		str = "quarter_max"
	} else {
		panic("ObservationValueType toField 오류:")
	}
	return str
}

const (
	//주
	WeekMinValue ObservationValueType = iota
	WeekMaxValue
	//월
	MonthMinValue
	MonthMaxValue
	//분기
	QuarterMinValue
	QuarterMaxValue
)

var ObservationValueTypeArr = []ObservationValueType{
	WeekMinValue, WeekMaxValue,
	MonthMinValue, MonthMaxValue,
	QuarterMinValue, QuarterMaxValue,
}
