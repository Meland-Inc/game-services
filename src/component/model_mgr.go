package component

import "github.com/Meland-Inc/game-services/src/application"

var mgrInstance *ModelManager

func GetInstance() *ModelManager {
	return mgrInstance
}

type ModelManager struct {
	appInter application.AppInterface
	models   map[string]ModelInterface
}

func InitModelManager(app application.AppInterface) *ModelManager {
	mgrInstance = &ModelManager{
		appInter: app,
		models:   make(map[string]ModelInterface),
	}
	return mgrInstance
}

func (mgr *ModelManager) GetApplication() application.AppInterface {
	return mgr.appInter
}

func (mgr *ModelManager) GetModel(modelName string) (ModelInterface, bool) {
	m, exist := mgr.models[modelName]
	return m, exist
}

func (mgr *ModelManager) AddModel(model ModelInterface) error {
	if err := model.OnInit(mgr); err != nil {
		return err
	}
	mgr.models[model.Name()] = model
	return nil
}

func (mgr *ModelManager) StartModel() error {
	for _, m := range mgr.models {
		if err := m.OnStart(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *ModelManager) TickModel(curMs int64) {
	for _, m := range mgr.models {
		m.OnTick(curMs)
	}
}

func (mgr *ModelManager) StopModel() error {
	for _, m := range mgr.models {
		if err := m.OnStop(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *ModelManager) ExitModel() error {
	for _, m := range mgr.models {
		if err := m.OnExit(); err != nil {
			return err
		}
	}
	return nil
}
