package module

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
)

var mgrInstance *ModuleManager

func GetInstance() *ModuleManager {
	return mgrInstance
}

func GetModel(moduleName string) (contract.IModuleInterface, bool) {
	return mgrInstance.GetModel(moduleName)
}

type ModuleManager struct {
	models map[string]contract.IModuleInterface
}

func InitModelManager() *ModuleManager {
	mgrInstance = &ModuleManager{
		models: make(map[string]contract.IModuleInterface),
	}
	return mgrInstance
}

func (mgr *ModuleManager) GetModel(moduleName string) (contract.IModuleInterface, bool) {
	m, exist := mgr.models[moduleName]
	return m, exist
}

func (mgr *ModuleManager) AddModel(model contract.IModuleInterface) error {
	serviceLog.Info("addModel [%v]", model.Name())
	if err := model.OnInit(); err != nil {
		return err
	}
	mgr.models[model.Name()] = model
	return nil
}

func (mgr *ModuleManager) StartModel() error {
	for _, m := range mgr.models {
		if err := m.OnStart(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *ModuleManager) TickModel(utc time.Time) {
	for _, m := range mgr.models {
		defer func() {
			if err := recover(); err != nil {
				serviceLog.Error("model [%v] panic: %+v", m.Name(), err)
			}
		}()
		m.OnTick(utc)
	}
}

func (mgr *ModuleManager) StopModel() error {
	for _, m := range mgr.models {
		if err := m.OnStop(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *ModuleManager) ExitModel() error {
	for _, m := range mgr.models {
		if err := m.OnExit(); err != nil {
			return err
		}
	}
	return nil
}
