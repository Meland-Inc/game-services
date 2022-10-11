package matrix

import (
	"math"
)

// //////////////////--LimitInt16--/////////////////////////
func LimitInt16(value int16, min int16, max int16) int16 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func LimitInt16Min(value int16, min int16) int16 {
	if value < min {
		return min
	}
	return value
}

func LimitInt16Max(value int16, max int16) int16 {
	if value > max {
		return max
	}
	return value
}

// //////////////////--LimitInt32--/////////////////////////
func LimitInt32(value int32, min int32, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func LimitInt32Min(value int32, min int32) int32 {
	if value < min {
		return min
	}
	return value
}

func LimitInt32Max(value int32, max int32) int32 {
	if value > max {
		return max
	}
	return value
}

// //////////////////--LimitInt--/////////////////////////
func LimitInt(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func LimitIntMin(value int, min int) int {
	if value < min {
		return min
	}
	return value
}

func LimitIntMax(value int, max int) int {
	if value > max {
		return max
	}
	return value
}

// //////////////////--LimitInt64--/////////////////////////
func LimitInt64(value int64, min int64, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func LimitInt64Min(value int64, min int64) int64 {
	if value < min {
		return min
	}
	return value
}

func LimitInt64Max(value int64, max int64) int64 {
	if value > max {
		return max
	}
	return value
}

// //////////////////--LimitFloat64--/////////////////////////
func LimitFloat64(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func LimitFloat64Min(value, min float64) float64 {
	if value < min {
		return min
	}
	return value
}

func LimitFloat64Max(value, max float64) float64 {
	if value > max {
		return max
	}
	return value
}

func AbsI16(val int16) int16 {
	if val > 0 {
		return val
	}
	return -val
}

func AbsI32(val int32) int32 {
	if val > 0 {
		return val
	}
	return -val
}

func AbsFloat64(val float64) float64 {
	if val >= 0.0 {
		return val
	}
	return -val
}

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func MaxInt16(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func MinInt16(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}

func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func Ceil(v float64) int32 {
	return int32(math.Ceil(v))
}

func Floor(v float64) int32 {
	return int32(math.Floor(v))
}

func RoundInt32(v float64) int32 {
	return int32(math.Floor(v + 0.5))
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
