package time_helper

import (
	"testing"
	"time"
)

func Test_times(t *testing.T) {
	var timeOffset int64 = 100

	tm := time.Now()
	t.Log(tm)
	tma := tm.Add(time.Duration(timeOffset) * time.Millisecond)
	t.Log(tma)
	t.Log(tm.UTC())
}
