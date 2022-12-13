package login_model

import (
	"fmt"
	"time"

	"github.com/Meland-Inc/game-services/src/global/component"
)

type LoginModel struct {
	component.ModelBase
}

func GetLoginModel() (*LoginModel, error) {
	iLoginModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_LOGIN)
	if !exist {
		return nil, fmt.Errorf("login model not found")
	}
	LoginModel, _ := iLoginModel.(*LoginModel)
	return LoginModel, nil
}

func NewLoginModel() *LoginModel {
	p := &LoginModel{}
	p.InitBaseModel(p, component.MODEL_NAME_LOGIN)
	return p
}

func (p *LoginModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("login model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)

	return nil
}

func (p *LoginModel) OnTick(utc time.Time) {
	p.ModelBase.OnTick(utc)
}
