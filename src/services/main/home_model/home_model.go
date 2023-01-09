package home_model

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/global/module"
)

type HomeModel struct {
	module.ModuleBase
}

func GetHomeModel() (*HomeModel, error) {
	iCtrlModel, exist := module.GetModel(module.MODULE_NAME_HOME)
	if !exist {
		return nil, fmt.Errorf("login model not found")
	}
	ctrlModel, _ := iCtrlModel.(*HomeModel)
	return ctrlModel, nil
}

func NewHomeModel() *HomeModel {
	p := &HomeModel{}
	p.InitBaseModel(p, module.MODULE_NAME_HOME)
	return p
}

func (p *HomeModel) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *HomeModel) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
}
