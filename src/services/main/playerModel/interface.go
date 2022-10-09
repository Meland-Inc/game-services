package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/configData"
)

func GetPlayerDataModel() (*PlayerDataModel, error) {
	iPlayerModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_PLAYER_DATA)
	if !exist {
		return nil, fmt.Errorf("player data model not found")
	}
	dataModel, _ := iPlayerModel.(*PlayerDataModel)
	return dataModel, nil
}

func (p *PlayerDataModel) OnPlayerDeath(
	userId int64,
	pos *proto.Vector3,
	killerId int64, killType proto.EntityType, KillerName string,
) error {
	sceneData, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return err
	}

	lvSetting := configData.ConfigMgr().RoleLevelCnf(sceneData.Level)
	if lvSetting == nil {
		return fmt.Errorf("role lv[%d] config not found", sceneData.Level)
	}
	if lvSetting.DeathExpLoss < 1 {
		return nil
	}

	// drop player current exp
	deathLossExp := lvSetting.DeathExpLoss
	if sceneData.Exp > deathLossExp {
		return p.DeductExp(userId, deathLossExp)
	}

	//  player max level item socket  lv -1
	maxLvSocket, err := p.playerMaxLevelItemSlot(userId)
	if err == nil && maxLvSocket.Level > 1 {
		_, err = p.setPlayerItemSlotLevel(userId, maxLvSocket.Position, maxLvSocket.Level-1, true)
	}
	return err
}
