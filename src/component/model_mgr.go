package component

import "github.com/Meland-Inc/game-services/src/application"

var mgrInstance *ModelManager

func GetInstance() *ModelManager {
	return mgrInstance
}

type ModelManager struct {
	appInter application.AppInterface
	models   map[string]Model
}

func InitModelManager(app application.AppInterface) *ModelManager {
	mgrInstance = &ModelManager{
		appInter: app,
		models:   make(map[string]Model),
	}
	return mgrInstance
}

func (mgr *ModelManager) GetApplication() application.AppInterface {
	return mgr.appInter
}

func (mgr *ModelManager) GetModel(modelName string) (Model, bool) {
	m, exist := mgr.models[modelName]
	return m, exist
}

func (mgr *ModelManager) AddModel(model Model) error {
	mgr.models[model.Name()] = model
	return model.OnInit(mgr)
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
