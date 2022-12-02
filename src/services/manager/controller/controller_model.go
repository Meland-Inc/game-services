package controller

import (
	"fmt"
	"sync"

	"github.com/Meland-Inc/game-services/src/global/component"
)

type ControllerModel struct {
	modelMgr   *component.ModelManager
	modelName  string
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
	model := &ControllerModel{}
	model.modelEvent = component.NewModelEvent(model)
	return model
}

func (p *ControllerModel) Name() string {
	return p.modelName
}

func (p *ControllerModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *ControllerModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("Controller model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_SERVICE_CONTROLLER
	return nil
}

func (p *ControllerModel) OnStart() (err error) {
	return nil
}

func (p *ControllerModel) OnTick(curMs int64) error {
	p.modelEvent.ReadEvent(curMs)
	return nil
}

func (p *ControllerModel) OnStop() error {
	p.modelMgr = nil
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
