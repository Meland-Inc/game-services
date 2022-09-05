package common

import "fmt"

func Assert(must bool, format string, v ...interface{}) {
	if !must {
		panic(fmt.Sprintf(format, v...))
	}
}
