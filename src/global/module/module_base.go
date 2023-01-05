package module

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/contract"
)

type ModuleBase struct {
	subModel  contract.IModuleInterface
	modelName string

	preTickSecond int
}

func (p *ModuleBase) Name() string { return p.modelName }

func (p *ModuleBase) InitBaseModel(subModel contract.IModuleInterface, modelName string) error {
	p.subModel = subModel
	p.modelName = modelName
	p.preTickSecond = time_helper.NowUTC().Second()
	return nil
}

func (p *ModuleBase) OnInit() error {
	return nil
}

func (p *ModuleBase) OnStart() (err error) {
	return nil
}

func (p *ModuleBase) OnTick(utc time.Time) {
	p.timeTick(utc)
}

func (p *ModuleBase) OnStop() error {
	return nil
}

func (p *ModuleBase) OnExit() error {
	return nil
}

func (p *ModuleBase) EventCall(env *ModuleEventReq) *ModuleEventResult {
	return nil
}

func (p *ModuleBase) EventCallNoReturn(env *ModuleEventReq) {}

func (p *ModuleBase) OnEvent(env *ModuleEventReq, curMs int64) {}

func (p *ModuleBase) timeTick(utc time.Time) {
	curSecond := utc.Second()
	if curSecond == p.preTickSecond {
		return
	}

	if curSecond == 0 {
		p.subModel.Minutely(utc)
	}
	if utc.Minute() == 0 && curSecond == 0 {
		p.subModel.Hourly(utc)
	}
	if utc.Hour() == 0 && utc.Minute() == 0 && curSecond == 0 {
		p.subModel.Daily(utc)
	}
	p.preTickSecond = curSecond
}

func (p *ModuleBase) Secondly(utc time.Time) {}

func (p *ModuleBase) Minutely(utc time.Time) {}

func (p *ModuleBase) Hourly(utc time.Time) {}

func (p *ModuleBase) Daily(utc time.Time) {}
