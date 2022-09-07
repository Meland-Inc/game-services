package time_helper

import (
	"time"
)

type TimeTemplate string

const (
	TimeTemplate1 TimeTemplate = "2006-01-02 15:04:05.000" // 常规类型
	TimeTemplate2 TimeTemplate = "2006/01/02 15:04:05.000" // 其他类型
	TimeTemplate3 TimeTemplate = "2006-01-02"              // 其他类型
	TimeTemplate4 TimeTemplate = "15:04:05"                // 其他类型
)

// Beginning of Day
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// End of Day
func Eod(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

func UnixMicro(t time.Time) int64 {
	return t.UnixNano() / 1e3

}

func UnixMill(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func UnixSec(t time.Time) int64 {
	return t.Unix()
}

func UnixMin(t time.Time) int64 {
	return t.Unix() / 60
}

func UnixHour(t time.Time) int64 {
	return t.Unix() / 3600
}

// 获取当前时间 精确到微秒
func NowMicro() int64 {
	return time.Now().UnixNano() / 1e3
}

// 获取当前时间 精确到毫秒
func NowMill() int64 {
	return time.Now().UnixNano() / 1e6
}

func NowSec() int64 {
	return UnixSec(time.Now())
}

// =========== CST 东八区时间================
func NowCST() time.Time {
	time.Local = time.FixedZone("CST", 3600*8)
	return time.Now().Local()
}

// 获取CST当前时间 精确到微秒
func NowCSTMicro() int64 {
	return NowCST().UnixNano() / 1e3
}

// 获取CST当前时间 精确到毫秒
func NowCSTMill() int64 {
	return NowCST().UnixNano() / 1e6
}

// 获取CST当前时间 精确到秒
func NowCSTSec() int64 {
	return UnixSec(NowCST())
}

/// =========================================

// =========== UTC 时间================
func NowUTC() time.Time {
	return time.Now().UTC()
}
func NowUTCMicro() int64 {
	return NowUTC().UnixNano() / 1e3
}
func NowUTCMill() int64 {
	return NowUTC().UnixNano() / 1e6
}
func NowUTCSec() int64 {
	return NowUTC().Unix()
}
func NowUTCMin(t time.Time) int64 {
	return NowUTCSec() / 60
}
func NowUTCHour(t time.Time) int64 {
	return NowUTCSec() / 3600
}

/// =========================================

func Format(sec int64, nsec int64, timeTemplate TimeTemplate) string {
	return time.Unix(sec, nsec).Format(string(timeTemplate))
}
