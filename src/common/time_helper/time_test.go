package time_helper

import (
	"testing"
	"time"
)

func Test_times(t *testing.T) {
	t.Log(time.Now())
	t.Log(NowCST())
}
