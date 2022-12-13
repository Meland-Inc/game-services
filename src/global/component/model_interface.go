package component

import "time"

type ModelInterface interface {
	Name() string
	OnInit(modelMgr *ModelManager) error
	OnStart() error
	OnTick(utc time.Time)
	OnStop() error
	OnExit() error
	EventCall(env *ModelEventReq) *ModelEventResult
	EventCallNoReturn(env *ModelEventReq)
	OnEvent(env *ModelEventReq, curMs int64)
	Secondly(utc time.Time)
	Minutely(utc time.Time)
	Hourly(utc time.Time)
	Daily(utc time.Time)
}
