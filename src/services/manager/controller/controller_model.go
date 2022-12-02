package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/global/component"
)

type ControllerModel struct {
	component.ModelBase
	modelEvent *component.ModelEvent

	controller         sync.Map
	startingPrivateSer sync.Map // { ownerId = *ServiceData} home service and Dungeon service
}

func GetControllerModel() (*ControllerModel, error) {
	iCtrlModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_SERVICE_CONTROLLER)
	if !exist {
		return nil, fmt.Errorf("login model not found")
	}
	ctrlModel, _ := iCtrlModel.(*ControllerModel)
	return ctrlModel, nil
}

func NewControllerModel() *ControllerModel {
	p := &ControllerModel{}
	p.modelEvent = component.NewModelEvent(p)
	p.InitBaseModel(p, component.MODEL_NAME_SERVICE_CONTROLLER)
	return p
}

func (p *ControllerModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("Controller model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)

	return nil
}

func (p *ControllerModel) OnStart() (err error) {
	return nil
}

func (p *ControllerModel) OnTick(utc time.Time) {
	p.ModelBase.OnTick(utc)
	p.modelEvent.ReadEvent(utc.UnixMilli())
}

func (p *ControllerModel) OnStop() error {
	p.ModelBase.OnStop()
	return nil
}

func (p *ControllerModel) OnExit() error {
	return nil
}

func (p *ControllerModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return p.modelEvent.EventCall(env)
}

func (p *ControllerModel) EventCallNoReturn(env *component.ModelEventReq) {
	p.modelEvent.EventCallNoReturn(env)
}

func (p *ControllerModel) Secondly(utc time.Time) {}
func (p *ControllerModel) Minutely(utc time.Time) {}
func (p *ControllerModel) Hourly(utc time.Time)   {}
func (p *ControllerModel) Daily(utc time.Time)    {}
