package userAgent

import (
	"fmt"
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/component"
)

type UserAgentModel struct {
	component.ModelBase
	record sync.Map
}

func NewUserAgentModel() *UserAgentModel {
	p := &UserAgentModel{}
	p.InitBaseModel(p, component.MODEL_NAME_USER_AGENT)
	return p
}

func GetUserAgentModel() *UserAgentModel {
	iUserAgentModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_USER_AGENT)
	if !exist {
		return nil
	}
	agentModel := iUserAgentModel.(*UserAgentModel)
	return agentModel
}

func (p *UserAgentModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)
	return nil
}

func (p *UserAgentModel) OnStop() error {
	p.ModelBase.OnStop()
	p.record = sync.Map{}
	return nil
}

func (p *UserAgentModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *UserAgentModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *UserAgentModel) OnEvent(env *component.ModelEventReq, curMs int64) {}

func (p *UserAgentModel) Secondly(utc time.Time) {}
func (p *UserAgentModel) Minutely(utc time.Time) {}
func (p *UserAgentModel) Hourly(utc time.Time)   {}
func (p *UserAgentModel) Daily(utc time.Time)    {}

func (p *UserAgentModel) GetUserAgent(userId int64) (*UserAgentData, bool) {
	iAgent, exist := p.record.Load(userId)
	if !exist {
		return nil, false
	}
	return iAgent.(*UserAgentData), exist
}

func (p *UserAgentModel) AllUserAgent() []*UserAgentData {
	agents := make([]*UserAgentData, 0)
	p.record.Range(func(key, value interface{}) bool {
		agents = append(agents, value.(*UserAgentData))
		return true
	})
	return agents
}

func (p *UserAgentModel) AllOnlineUserIds() []int64 {
	userIds := make([]int64, 0)
	p.record.Range(func(key, value interface{}) bool {
		userIds = append(userIds, value.(*UserAgentData).UserId)
		return true
	})
	return userIds
}

func (p *UserAgentModel) AddUserAgentRecord(
	userId int64,
	agentAppId, socketId, sceneAppId string,
) (*UserAgentData, error) {
	if userId == 0 || agentAppId == "" || socketId == "" {
		return nil, fmt.Errorf("user agent data is invalid")
	}

	agentData := &UserAgentData{
		AgentAppId:          agentAppId,
		SocketId:            socketId,
		InSceneServiceAppId: sceneAppId,
		UserId:              userId,
		LoginAt:             time_helper.NowUTCMill(),
	}
	p.record.Store(userId, agentData)
	return agentData, nil
}

func (p *UserAgentModel) CheckAndAddUserAgentRecord(
	userId int64,
	agentAppId, socketId, sceneAppId string,
) (*UserAgentData, error) {
	return p.AddUserAgentRecord(userId, agentAppId, socketId, sceneAppId)
}

func (p *UserAgentModel) RemoveUserAgentRecord(userId int64) {
	p.record.Delete(userId)
}
