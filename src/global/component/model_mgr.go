package component

import (
	"time"

	"github.com/Meland-Inc/game-services/src/application"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

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
	serviceLog.Info("addModel [%v]", model.Name())
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

func (mgr *ModelManager) TickModel(utc time.Time) {
	for _, m := range mgr.models {
		defer func() {
			if err := recover(); err != nil {
				serviceLog.Error("model [%v] panic: %+v", m.Name(), err)
			}
		}()
		m.OnTick(utc)
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
