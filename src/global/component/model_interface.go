package component

type ModelInterface interface {
	Name() string
	ModelMgr() *ModelManager
	OnInit(modelMgr *ModelManager) error
	OnStart() error
	OnTick(curMs int64) error
	OnStop() error
	OnExit() error
}
