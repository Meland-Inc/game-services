package module

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/dapr/go-sdk/service/common"
)

type ModuleEventResult struct {
	result interface{}
	err    error
}

func NewModuleEventResult(result interface{}, err error) *ModuleEventResult {
	return &ModuleEventResult{
		result: result,
		err:    err,
	}
}

func (r *ModuleEventResult) GetError() error {
	return r.err
}
func (r *ModuleEventResult) SetError(err error) {
	r.err = err
}
func (r *ModuleEventResult) GetResult() interface{} {
	return r.result
}
func (r *ModuleEventResult) SetResult(result interface{}) {
	r.result = result
}

type ModuleEventReq struct {
	eventType  string
	msg        interface{}
	mustReturn bool
	result     chan contract.IModuleEventResult
}

func NewModuleEventReq(
	eventType string,
	msg interface{},
	mustReturn bool,
	result chan contract.IModuleEventResult,
) *ModuleEventReq {
	return &ModuleEventReq{
		eventType:  eventType,
		msg:        msg,
		mustReturn: mustReturn,
		result:     result,
	}
}
func (e *ModuleEventReq) SetEventType(eventType string) {
	e.eventType = eventType
}
func (e *ModuleEventReq) GetEventType() string {
	return e.eventType
}

func (e *ModuleEventReq) SetMsg(msg interface{}) {
	e.msg = msg
}
func (e *ModuleEventReq) GetMsg() interface{} {
	return e.msg
}

func (e *ModuleEventReq) SetResultChan(result chan contract.IModuleEventResult) {
	e.result = result
}
func (e *ModuleEventReq) GetResultChan() chan contract.IModuleEventResult {
	return e.result
}

func (e *ModuleEventReq) SetMustReturn(bSet bool) { e.mustReturn = bSet }
func (e *ModuleEventReq) GetMustReturn() bool     { return e.mustReturn }

func (e *ModuleEventReq) WriteResult(res contract.IModuleEventResult) {
	if !e.mustReturn {
		return
	}
	if nil == e.result {
		serviceLog.Error("WriteResult fail. result is nil. msg(%v)", e)
		return
	}

	// 写超时检测(用于检测有多个向同一个channel中写数据的非法操作)
	select {
	case <-time.After(time.Second):
		serviceLog.Error("WriteResult fail. write overtime. msg(%v)", e)
		return
	case e.result <- res:
		return
	}
}

func (e *ModuleEventReq) UnmarshalToDaprEventData(v interface{}) error {
	msg, ok := e.GetMsg().(*common.TopicEvent)
	if !ok {
		return fmt.Errorf("Unmarshal dapr Event failed: %v", msg)
	}
	return grpcNetTool.UnmarshalGrpcTopicEvent(msg, v)
}

func (e *ModuleEventReq) UnmarshalToDaprCallData(v interface{}) error {
	inputBs, ok := e.GetMsg().([]byte)
	if !ok {
		return fmt.Errorf("Unmarshal dapr call failed: %s", inputBs)

	}
	return grpcNetTool.UnmarshalGrpcData(inputBs, v)
}

type ModuleEvent struct {
	eventChan chan contract.IModuleEventReq
}

func NewModelEvent() *ModuleEvent {
	return &ModuleEvent{
		eventChan: make(chan contract.IModuleEventReq, 3000),
	}
}

func (p *ModuleEvent) eventCallImpl(env contract.IModuleEventReq, mustReturn bool) contract.IModuleEventResult {
	if len(p.eventChan) > 200 {
		serviceLog.Warning("ModuleEvent Msg length (%d)>200", len(p.eventChan))
	}

	env.SetMustReturn(mustReturn)
	var outCh chan contract.IModuleEventResult = nil
	if env.GetMustReturn() {
		env.SetResultChan(make(chan contract.IModuleEventResult, 1))
		outCh = env.GetResultChan()
	}

	p.eventChan <- env

	if env.GetMustReturn() {
		select {
		case <-time.After(time.Second * 5):
			result := &ModuleEventResult{
				err: fmt.Errorf("ModuleEvent timeout. msgType(%v),  msg: %+v", env.GetEventType(), env.GetMsg()),
			}
			serviceLog.Error(result.err.Error())
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

func (p *ModuleEvent) EventCall(env contract.IModuleEventReq) contract.IModuleEventResult {
	return p.eventCallImpl(env, true)
}

func (p *ModuleEvent) EventCallNoReturn(env contract.IModuleEventReq) {
	p.eventCallImpl(env, false)
}

func (p *ModuleEvent) ReadEvent() contract.IModuleEventReq {
	select {
	case e := <-p.eventChan:
		return e
	default:
		return nil
	}
}
