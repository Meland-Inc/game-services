package logicDemoModel

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/global/module"
)

type LogicDemoModel struct {
	module.ModuleBase
}

func GetLogicDemoModel() (*LogicDemoModel, error) {
	iLogicDemoModel, exist := module.GetModel(module.MODULE_NAME_LOGIC_DEMO)
	if !exist {
		return nil, fmt.Errorf("logic demo model not found")
	}
	LogicDemoModel, _ := iLogicDemoModel.(*LogicDemoModel)
	return LogicDemoModel, nil
}

func NewLogicDemoModel() *LogicDemoModel {
	p := &LogicDemoModel{}
	p.InitBaseModel(p, module.MODULE_NAME_LOGIC_DEMO)
	return p
}

func (p *LogicDemoModel) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *LogicDemoModel) OnStart() (err error) {
	p.ModuleBase.OnStart()
	return err
}

func (p *LogicDemoModel) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
}

func (p *LogicDemoModel) Secondly(utc time.Time) {}
func (p *LogicDemoModel) Minutely(utc time.Time) {}
func (p *LogicDemoModel) Hourly(utc time.Time)   {}
func (p *LogicDemoModel) Daily(utc time.Time)    {}
