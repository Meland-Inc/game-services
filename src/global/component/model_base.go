package component

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

type ModelBase struct {
	modelMgr  *ModelManager
	subModel  ModelInterface
	modelName string

	preTickSecond int
}

func (p *ModelBase) Name() string            { return p.modelName }
func (p *ModelBase) ModelMgr() *ModelManager { return p.modelMgr }

func (p *ModelBase) InitBaseModel(subModel ModelInterface, modelName string) error {
	p.subModel = subModel
	p.modelName = modelName
	p.preTickSecond = time_helper.NowUTC().Second()
	p.modelMgr = GetInstance()
	return nil
}

func (p *ModelBase) OnInit(modelMgr *ModelManager) error {
	p.modelMgr = modelMgr
	return nil
}

func (p *ModelBase) OnStart() (err error) {
	return nil
}

func (p *ModelBase) OnTick(utc time.Time) {
	p.timeTick(utc)
}

func (p *ModelBase) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *ModelBase) OnExit() error {
	return nil
}

func (p *ModelBase) EventCall(env *ModelEventReq) *ModelEventResult {
	return nil
}

func (p *ModelBase) EventCallNoReturn(env *ModelEventReq) {}

func (p *ModelBase) OnEvent(env *ModelEventReq, curMs int64) {}

func (p *ModelBase) timeTick(utc time.Time) {
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

func (p *ModelBase) Secondly(utc time.Time) {}

func (p *ModelBase) Minutely(utc time.Time) {}

func (p *ModelBase) Hourly(utc time.Time) {}

func (p *ModelBase) Daily(utc time.Time) {}
