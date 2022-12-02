package login_model

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/global/component"
)

type LoginModel struct {
	modelMgr  *component.ModelManager
	modelName string
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
	return &LoginModel{}
}

func (p *LoginModel) Name() string {
	return p.modelName
}

func (p *LoginModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *LoginModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("login model init service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_LOGIN
	return nil
}

func (p *LoginModel) OnStart() (err error) {
	return nil
}

func (p *LoginModel) OnTick(curMs int64) error {
	return nil
}

func (p *LoginModel) OnStop() error {
	p.modelMgr = nil
	return nil
}

func (p *LoginModel) OnExit() error {
	return nil
}

func (p *LoginModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *LoginModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *LoginModel) OnEvent(env *component.ModelEventReq, curMs int64) {}
