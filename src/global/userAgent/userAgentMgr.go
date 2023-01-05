package userAgent

import (
	"fmt"
	"game-message-core/proto"
	"sync"
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"

	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
)

type UserAgentModel struct {
	module.ModuleBase
	record sync.Map
}

func NewUserAgentModel() *UserAgentModel {
	p := &UserAgentModel{}
	p.InitBaseModel(p, module.MODULE_NAME_USER_AGENT)
	return p
}

func GetUserAgentModel() *UserAgentModel {
	iUserAgentModel, exist := module.GetModel(module.MODULE_NAME_USER_AGENT)
	if !exist {
		return nil
	}
	agentModel := iUserAgentModel.(*UserAgentModel)
	return agentModel
}

func (p *UserAgentModel) OnInit() error {
	p.ModuleBase.OnInit()
	return nil
}

func (p *UserAgentModel) OnStop() error {
	p.ModuleBase.OnStop()
	p.record = sync.Map{}
	return nil
}

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

func SendToPlayer(serviceAppId string, userId int64, msg *proto.Envelope) error {
	agent, exist := GetUserAgentModel().GetUserAgent(userId)
	if !exist {
		return fmt.Errorf("user [%d] agent data not found", userId)
	}
	return agent.SendToPlayer(serviceAppId, msg)
}

func (p *UserAgentModel) EventCall(env contract.IModuleEventReq) contract.IModuleEventResult {
	return nil
}
func (p *UserAgentModel) EventCallNoReturn(env contract.IModuleEventReq) {}
func (p *UserAgentModel) ReadEvent() contract.IModuleEventReq {
	return nil
}
