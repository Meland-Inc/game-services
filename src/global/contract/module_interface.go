package contract

import "time"

type IModuleInterface interface {
	Name() string
	OnInit() error
	OnStart() error
	OnTick(utc time.Time)
	OnStop() error
	OnExit() error
	Secondly(utc time.Time)
	Minutely(utc time.Time)
	Hourly(utc time.Time)
	Daily(utc time.Time)
}

type IModuleEventReq interface {
	SetEventType(eventType string)
	GetEventType() string
	SetMsg(msg interface{})
	GetMsg() interface{}
	SetResultChan(result chan IModuleEventResult)
	GetResultChan() chan IModuleEventResult
	SetMustReturn(bSet bool)
	GetMustReturn() bool
	WriteResult(ret IModuleEventResult)
	UnmarshalToDaprEventData(v interface{}) error
	UnmarshalToDaprCallData(v interface{}) error
}

type IModuleEventResult interface {
	GetError() error
	SetError(err error)
	GetResult() interface{}
	SetResult(result interface{})
}

type IModuleEvent interface {
	EventCall(env IModuleEventReq) IModuleEventResult
	EventCallNoReturn(env IModuleEventReq)
	ReadEvent() IModuleEventReq
}
