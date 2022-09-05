package serviceLog

import (
	"testing"
	"time"
)

func TestSetLogger(t *testing.T) {
	SetWriteFile("bin/logs/abdc.log", false)
	now := time.Now()
	s := 60*1e9 - (now.Second()*1e9 + now.Nanosecond())
	s2 := (60-now.Second())*1e9 - now.Nanosecond()
	Info("xxxxxxxxxxxxxxx %d %d %d %d", s, s2, now.Nanosecond(), s)
}
