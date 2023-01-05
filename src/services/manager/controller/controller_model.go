package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
)

type ControllerModel struct {
	module.ModuleBase
	modelEvent *module.ModuleEvent

	controller         sync.Map
	startingPrivateSer sync.Map // { ownerId = *ServiceData} home service and Dungeon service
}

func GetControllerModel() (*ControllerModel, error) {
	iCtrlModel, exist := module.GetModel(module.MODULE_NAME_SERVICE_CONTROLLER)
	if !exist {
		return nil, fmt.Errorf("login model not found")
	}
	ctrlModel, _ := iCtrlModel.(*ControllerModel)
	return ctrlModel, nil
}

func NewControllerModel() *ControllerModel {
	p := &ControllerModel{}
	p.modelEvent = module.NewModelEvent()
	p.InitBaseModel(p, module.MODULE_NAME_SERVICE_CONTROLLER)
	return p
}

func (p *ControllerModel) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *ControllerModel) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
	if env := p.ReadEvent(); env != nil {
		p.OnEvent(env, utc.UnixMilli())
	}
}

func (p *ControllerModel) OnStop() error {
	p.ModuleBase.OnStop()
	return nil
}

func (p *ControllerModel) Secondly(utc time.Time) {
	p.checkAndRemoveTimeOutSer(utc.UnixMilli())
}
func (p *ControllerModel) Minutely(utc time.Time) {}
func (p *ControllerModel) Hourly(utc time.Time)   {}
func (p *ControllerModel) Daily(utc time.Time)    {}

func (p *ControllerModel) EventCall(env contract.IModuleEventReq) contract.IModuleEventResult {
	return p.modelEvent.EventCall(env)
}

func (p *ControllerModel) EventCallNoReturn(env contract.IModuleEventReq) {
	p.modelEvent.EventCallNoReturn(env)
}

func (p *ControllerModel) ReadEvent() contract.IModuleEventReq {
	return p.modelEvent.ReadEvent()
}
