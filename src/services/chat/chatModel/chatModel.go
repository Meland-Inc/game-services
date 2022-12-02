package chatModel

import (
	"fmt"

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
	modelMgr  *component.ModelManager
	modelName string
	mapGrids  map[int32]*MapGrid
	Players   map[int64]*PlayerChatData
}

func NewChatModel() *ChatModel {
	return &ChatModel{
		mapGrids: make(map[int32]*MapGrid),
		Players:  make(map[int64]*PlayerChatData),
	}
}

func (p *ChatModel) Name() string {
	return p.modelName
}

func (p *ChatModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *ChatModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("chat model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_CHAT
	return nil
}

func (p *ChatModel) OnStart() error {
	p.onStart()
	return nil
}

func (p *ChatModel) OnTick(curMs int64) error {
	return nil
}

func (p *ChatModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *ChatModel) OnExit() error {
	return nil
}

func (p *ChatModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *ChatModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *ChatModel) OnEvent(env *component.ModelEventReq, curMs int64) {}
