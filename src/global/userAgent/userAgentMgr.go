package userAgent

import (
	"fmt"
	"sync"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type UserAgentModel struct {
	modelMgr  *component.ModelManager
	modelName string
	record    sync.Map
}

func NewUserAgentModel() *UserAgentModel {
	return &UserAgentModel{}
}

func GetUserAgentModel() *UserAgentModel {
	iUserAgentModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_USER_AGENT)
	if !exist {
		return nil
	}
	agentModel := iUserAgentModel.(*UserAgentModel)
	return agentModel
}

func (p *UserAgentModel) Name() string {
	return p.modelName
}

func (p *UserAgentModel) ModelMgr() *component.ModelManager {
	return p.modelMgr
}

func (p *UserAgentModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("service model manager is nil")
	}
	p.modelMgr = modelMgr
	p.modelName = component.MODEL_NAME_USER_AGENT
	return nil
}

func (p *UserAgentModel) OnStart() error {
	return nil
}

func (p *UserAgentModel) OnTick(curMs int64) error {
	return nil
}

func (p *UserAgentModel) OnStop() error {
	p.record = sync.Map{}
	p.modelMgr = nil
	return nil
}

func (p *UserAgentModel) OnExit() error {
	return nil
}

func (p *UserAgentModel) GetUserAgent(userId int64) (*UserAgentData, bool) {
	iAgent, exist := p.record.Load(userId)
	if !exist {
		return nil, false
	}
	return iAgent.(*UserAgentData), exist
}

func (p *UserAgentModel) AddUserAgentRecord(userId int64, agentAppId, socketId string) error {
	if userId == 0 || agentAppId == "" || socketId == "" {
		return fmt.Errorf("user agent data is invalid")
	}

	agentData := &UserAgentData{
		AgentAppId: agentAppId,
		SocketId:   socketId,
		UserId:     userId,
		LoginAt:    time_helper.NowUTCMill(),
	}
	p.record.Store(userId, agentData)
	return nil
}

func (p *UserAgentModel) CheckAndAddUserAgentRecord(userId int64, agentAppId, socketId string) error {
	return p.AddUserAgentRecord(userId, agentAppId, socketId)
}

func (p *UserAgentModel) RemoveUserAgentRecord(userId int64) {
	p.record.Delete(userId)
}
