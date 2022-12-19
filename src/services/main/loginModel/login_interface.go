package login_model

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
	"gorm.io/gorm"
)

func (p *LoginModel) selectLoginData(userId int64) (*dbData.LoginData, error) {
	loginData := &dbData.LoginData{}
	err := gameDB.GetGameDB().Where("user_id = ?", userId).First(loginData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return loginData, nil
}

func (p *LoginModel) OnLogin(userId int64, agentAppId, socketId, sceneAppId string) error {
	loginData, err := p.selectLoginData(userId)
	if err != nil {
		return err
	}

	if loginData != nil {
		// event call tick out old login player
		grpcPubsubEvent.RPCPubsubEventTickOutPlayer(
			userId,
			loginData.AgentAppId,
			loginData.SocketId,
			loginData.InSceneService,
			proto.TickOutType_RepeatSingIn,
		)
	}

	loginData = dbData.NewLoginData(userId, agentAppId, socketId, sceneAppId)
	return gameDB.GetGameDB().Save(loginData).Error
}

func (p *LoginModel) OnLogOut(userId int64) {
	err := gameDB.GetGameDB().Where("user_id=?", userId).Delete(&dbData.LoginData{}).Error
	if err != nil {
		serviceLog.Error(err.Error())
	}
}

func (p *LoginModel) GetUserLoginData(userId int64, agentAppId, socketId string) (sceneAppId string, err error) {
	row := &dbData.PlayerSceneData{}
	err = gameDB.GetGameDB().Where("user_id = ?", userId).First(row).Error
	if err != nil {
		return sceneAppId, err
	}

	_, subType, exist := configData.ConfigMgr().GetSceneArea(row.MapId)
	if !exist || subType == proto.SceneServiceSubType_Home {
		enterMap, enterPos := configData.GetHomeEntrance()
		if err = gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
			row.MapId = enterMap
			row.X = enterPos.X
			row.Y = enterPos.Y
			row.Z = enterPos.Z
			return tx.Save(row).Error
		}); err != nil {
			return sceneAppId, err
		}
	}

	serviceOut, err := grpcInvoke.RPCSelectService(
		proto.ServiceType_ServiceTypeScene,
		proto.SceneServiceSubType_World,
		0,
		row.MapId,
	)
	serviceLog.Debug("select world output = %+v, err:%+v", serviceOut, err)
	if err != nil {
		return "", err
	}
	if serviceOut.ErrorCode > 0 {
		return "", fmt.Errorf(serviceOut.ErrorMessage)
	}

	sceneAppId = serviceOut.Service.AppId
	err = p.OnLogin(userId, agentAppId, socketId, sceneAppId)
	return sceneAppId, err
}
