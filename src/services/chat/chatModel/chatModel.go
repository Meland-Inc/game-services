package chatModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
)

func GetChatModel() (*ChatModel, error) {
	iModel, exist := module.GetModel(module.MODULE_NAME_CHAT)
	if !exist {
		return nil, fmt.Errorf("chat view grid model not found")
	}
	model, _ := iModel.(*ChatModel)
	return model, nil
}

type ChatModel struct {
	module.ModuleBase
	mapGrids map[int32]*MapGrid
	Players  map[int64]*PlayerChatData
}

func NewChatModel() *ChatModel {
	p := &ChatModel{
		mapGrids: make(map[int32]*MapGrid),
		Players:  make(map[int64]*PlayerChatData),
	}
	p.InitBaseModel(p, module.MODULE_NAME_CHAT)
	return p
}

func (p *ChatModel) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *ChatModel) OnStart() error {
	p.ModuleBase.OnStart()
	p.onStart()
	return nil
}

func (p *ChatModel) OnTick(utc time.Time) {
}
func (p *ChatModel) Secondly(utc time.Time) {}
func (p *ChatModel) Minutely(utc time.Time) {}
func (p *ChatModel) Hourly(utc time.Time)   {}
func (p *ChatModel) Daily(utc time.Time)    {}

func (p *ChatModel) EventCall(env contract.IModuleEventReq) contract.IModuleEventResult {
	return nil
}
func (p *ChatModel) EventCallNoReturn(env contract.IModuleEventReq) {}
func (p *ChatModel) ReadEvent() contract.IModuleEventReq {
	return nil
}
