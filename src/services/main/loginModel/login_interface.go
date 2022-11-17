package login_model

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
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
