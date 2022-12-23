package time_helper

import (
	"testing"
	"time"
)

func Test_times(t *testing.T) {
	// var timeOffset int64 = 100

	// tm := time.Now().UTC()
	// t.Log(tm)
	// // t.Log(tm.Local())
	// ms := tm.UnixMilli()
	// t.Log(ms)

	// tma := tm.Add(time.Duration(timeOffset) * time.Millisecond)
	// t.Log(tma)
	// t.Log(tm.UTC())

	var ms int64 = 1671791513368
	tm2 := time.UnixMilli(ms)
	t.Log(tm2.UTC())
	t.Log(tm2)
	// t.Log(tm2.UTC().UnixMilli() == tm.UnixMilli())
}
