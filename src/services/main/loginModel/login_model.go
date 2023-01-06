package login_model

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/global/module"
)

type LoginModel struct {
	module.ModuleBase
}

func GetLoginModel() (*LoginModel, error) {
	iLoginModel, exist := module.GetModel(module.MODULE_NAME_LOGIN)
	if !exist {
		return nil, fmt.Errorf("login model not found")
	}
	LoginModel, _ := iLoginModel.(*LoginModel)
	return LoginModel, nil
}

func NewLoginModel() *LoginModel {
	p := &LoginModel{}
	p.InitBaseModel(p, module.MODULE_NAME_LOGIN)
	return p
}

func (p *LoginModel) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *LoginModel) OnTick(utc time.Time) {
	p.ModuleBase.OnTick(utc)
}
