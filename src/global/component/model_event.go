package component

import (
	"fmt"
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

type ModelEventReq struct {
	EventType  string
	Msg        interface{}
	mustReturn bool
	Result     chan *ModelEventResult
}

func (e *ModelEventReq) GetMsg() interface{} {
	return e.Msg
}
func (e *ModelEventReq) GetResultChan() chan *ModelEventResult {
	return e.Result
}
func (e *ModelEventReq) SetResultChan(result chan *ModelEventResult) {
	e.Result = result
}
func (e *ModelEventReq) MustReturn() bool {
	return e.mustReturn
}
func (e *ModelEventReq) SetMustReturn(bSet bool) {
	e.mustReturn = bSet
}

func (e *ModelEventReq) WriteResult(ret *ModelEventResult) {
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

type ModelEvent struct {
	model     ModelInterface
	eventChan chan *ModelEventReq
}

func NewModelEvent(model ModelInterface) *ModelEvent {
	return &ModelEvent{
		model:     model,
		eventChan: make(chan *ModelEventReq, 3000),
	}
}

func (p *ModelEvent) eventCallImpl(env *ModelEventReq, mustReturn bool) *ModelEventResult {
	if len(p.eventChan) > 200 {
		serviceLog.Warning("%s event lenMsg(%d)>200", p.model.Name(), len(p.eventChan))
	}

	env.SetMustReturn(mustReturn)
	var outCh chan *ModelEventResult = nil
	if env.MustReturn() {
		env.SetResultChan(make(chan *ModelEventResult, 1))
		outCh = env.GetResultChan()
	}

	p.eventChan <- env

	if env.MustReturn() {
		select {
		case <-time.After(time.Second * 5):
			result := &ModelEventResult{
				Err: fmt.Errorf("%s event timeout. msgType(%v),  msg: %+v", p.model.Name(), env.EventType, env.Msg),
			}
			serviceLog.Error(result.Err.Error())
			return result

		case retMsg := <-env.GetResultChan():
			// 数据回写给外部调用者的channel
			if nil != outCh {
				outCh <- retMsg
			}
			return retMsg
		}
	}
	return nil
}

func (p *ModelEvent) EventCall(env *ModelEventReq) *ModelEventResult {
	return p.eventCallImpl(env, true)
}

func (p *ModelEvent) EventCallNoReturn(env *ModelEventReq) {
	p.eventCallImpl(env, false)
}

func (p *ModelEvent) ReadEvent(curMs int64) {
	select {
	case e := <-p.eventChan:
		p.model.OnEvent(e, curMs)
	default:
		break
	}
}
