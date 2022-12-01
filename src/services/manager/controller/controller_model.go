package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type ControllerModel struct {
	modelMgr  *component.ModelManager
	modelName string

	eventChan chan *component.ModelEvent

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
	return &ControllerModel{
		eventChan: make(chan *component.ModelEvent, 3000),
	}
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
	p.eventTick(curMs)
	return nil
}

func (p *ControllerModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *ControllerModel) OnExit() error {
	return nil
}

func (p *ControllerModel) eventCallImpl(env *component.ModelEvent, mustReturn bool) *component.ModelEventResult {
	if len(p.eventChan) > 200 {
		serviceLog.Warning("ControllerModel event lenMsg(%d)>200", len(p.eventChan))
	}

	env.SetMustReturn(mustReturn)
	var outCh chan *component.ModelEventResult = nil
	if env.MustReturn() {
		env.SetResultChan(make(chan *component.ModelEventResult, 1))
		outCh = env.GetResultChan()
	}

	p.eventChan <- env

	if env.MustReturn() {
		select {
		case <-time.After(time.Second * 5):
			result := &component.ModelEventResult{
				Err: fmt.Errorf("ControllerModel event timeout. msgType(%v),  msg: %+v", env.EventType, env.Msg),
			}
			serviceLog.Error(result.Err.Error())
			return result

		case retMsg := <-env.GetResultChan():
			// 数据回写给外部调用者的channel
			if nil != outCh {
				outCh <- retMsg
			}
			return retMsg
		}
	}
	return nil
}

func (p *ControllerModel) EventCall(env *component.ModelEvent) *component.ModelEventResult {
	return p.eventCallImpl(env, true)
}

func (p *ControllerModel) EventCallNoReturn(env *component.ModelEvent) {
	p.eventCallImpl(env, false)
}

func (p *ControllerModel) eventTick(curMs int64) {
	select {
	case e := <-p.eventChan:
		p.onEvent(e, curMs)
	default:
		break
	}
}
