package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// [min, max]
func Random(min int, max int) int {
	if min >= max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}

// [min, max]
func Random16(min int16, max int16) int16 {
	if min >= max {
		min, max = max, min
	}
	return int16(rand.Int31n(int32(max-min+1))) + min
}

// [min, max]
func Random32(min int32, max int32) int32 {
	if min >= max {
		min, max = max, min
	}
	return rand.Int31n(max-min+1) + min
}

// [min, max]
func RandomU32(min uint32, max uint32) uint32 {
	if min >= max {
		min, max = max, min
	}
	return uint32(rand.Int63n(int64(max)-int64(min)+1) + int64(min))
}

// [min, max]
func Random64(min int64, max int64) int64 {
	if min >= max {
		min, max = max, min
	}
	return rand.Int63n(max-min+1) + min
}
