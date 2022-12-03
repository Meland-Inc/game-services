package home_model

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/global/component"
)

type HomeModel struct {
	component.ModelBase
	modelEvent *component.ModelEvent
}

func GetHomeModel() (*HomeModel, error) {
	iCtrlModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_HOME)
	if !exist {
		return nil, fmt.Errorf("login model not found")
	}
	ctrlModel, _ := iCtrlModel.(*HomeModel)
	return ctrlModel, nil
}

func NewHomeModel() *HomeModel {
	p := &HomeModel{}
	p.modelEvent = component.NewModelEvent(p)
	p.InitBaseModel(p, component.MODEL_NAME_HOME)
	return p
}

func (p *HomeModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("Home model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)
	return nil
}

func (p *HomeModel) OnTick(utc time.Time) {
	p.modelEvent.ReadEvent(utc.UnixMilli())
}

func (p *HomeModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return p.modelEvent.EventCall(env)
}

func (p *HomeModel) EventCallNoReturn(env *component.ModelEventReq) {
	p.modelEvent.EventCallNoReturn(env)
}
