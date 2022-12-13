package chatModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/global/component"
)

func GetChatModel() (*ChatModel, error) {
	iModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_CHAT)
	if !exist {
		return nil, fmt.Errorf("chat view grid model not found")
	}
	model, _ := iModel.(*ChatModel)
	return model, nil
}

type ChatModel struct {
	component.ModelBase
	mapGrids map[int32]*MapGrid
	Players  map[int64]*PlayerChatData
}

func NewChatModel() *ChatModel {
	p := &ChatModel{
		mapGrids: make(map[int32]*MapGrid),
		Players:  make(map[int64]*PlayerChatData),
	}
	p.InitBaseModel(p, component.MODEL_NAME_CHAT)
	return p
}

func (p *ChatModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("chat model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)
	return nil
}

func (p *ChatModel) OnStart() error {
	p.ModelBase.OnStart()
	p.onStart()
	return nil
}

func (p *ChatModel) OnTick(utc time.Time) {
}

func (p *ChatModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *ChatModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *ChatModel) OnEvent(env *component.ModelEventReq, curMs int64) {}

func (p *ChatModel) Secondly(utc time.Time) {}
func (p *ChatModel) Minutely(utc time.Time) {}
func (p *ChatModel) Hourly(utc time.Time)   {}
func (p *ChatModel) Daily(utc time.Time)    {}
