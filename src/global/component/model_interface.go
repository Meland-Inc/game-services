package component

type ModelInterface interface {
	Name() string
	ModelMgr() *ModelManager
	OnInit(modelMgr *ModelManager) error
	OnStart() error
	OnTick(curMs int64) error
	OnStop() error
	OnExit() error
	EventCall(env *ModelEventReq) *ModelEventResult
	EventCallNoReturn(env *ModelEventReq)
	OnEvent(env *ModelEventReq, curMs int64)
}
