package component

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

type ModelEventResult struct {
	Result interface{}
	Err    error
}

func (r *ModelEventResult) GetError() error {
	return r.Err
}
func (r *ModelEventResult) SetError(err error) {
	r.Err = err
}
func (r *ModelEventResult) GetResult() interface{} {
	return r.Result
}
func (r *ModelEventResult) SetResult(result interface{}) {
	r.Result = result
}

type ModelEvent struct {
	EventType  string
	Msg        interface{}
	mustReturn bool
	Result     chan *ModelEventResult
}

func (e *ModelEvent) GetMsg() interface{} {
	return e.Msg
}
func (e *ModelEvent) GetResultChan() chan *ModelEventResult {
	return e.Result
}
func (e *ModelEvent) SetResultChan(result chan *ModelEventResult) {
	e.Result = result
}
func (e *ModelEvent) MustReturn() bool {
	return e.mustReturn
}
func (e *ModelEvent) SetMustReturn(bSet bool) {
	e.mustReturn = bSet
}

func (e *ModelEvent) WriteResult(ret *ModelEventResult) {
	if !e.mustReturn {
		return
	}
	if nil == e.Result {
		serviceLog.Error("model WriteResult fail. result is nil. msg(%v)", e)
		return
	}

	// 写超时检测(用于检测有多个向同一个channel中写数据的非法操作)
	select {
	case <-time.After(time.Second):
		serviceLog.Error("ZoneEvent:WriteResult fail. write overtime. msg(%v)", e)
		return
	case e.Result <- ret:
		return
	}
}
