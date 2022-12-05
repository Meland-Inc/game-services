package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

func (p *PlayerDataModel) GetPlayerBaseData(userId int64) (*dbData.PlayerBaseData, error) {
	baseData := &dbData.PlayerBaseData{}
	err := gameDB.GetGameDB().Where("user_id = ?", userId).First(baseData).Error
	return baseData, err
}

func (p *PlayerDataModel) PlayerAllData(userId int64) (
	baseData *dbData.PlayerBaseData,
	sceneData *dbData.PlayerSceneData,
	avatars []*Item,
	profile *proto.EntityProfile,
	err error,
) {
	if baseData, err = p.GetPlayerBaseData(userId); err != nil {
		return
	}
	if sceneData, err = p.GetPlayerSceneData(userId); err != nil {
		return
	}
	if avatars, err = p.UsingAvatars(userId); err != nil {
		return
	}
	profile, err = p.GetPlayerProfile(userId)
	return
}

func (p *PlayerDataModel) PlayerProtoData(userId int64) (*proto.Player, error) {
	baseData, sceneData, avatars, profile, err := p.PlayerAllData(userId)
	if err != nil {
		return nil, err
	}

	pos := &proto.Vector3{X: sceneData.X, Y: sceneData.Y, Z: sceneData.Z}
	dir := &proto.Vector3{X: sceneData.DirX, Y: sceneData.DirY, Z: sceneData.DirZ}
	player := &proto.Player{
		BaseData: baseData.ToNetPlayerBaseData(),
		Profile:  profile,
		Active:   sceneData.Hp > 0,
		MapId:    sceneData.MapId,
		Position: pos,
		Dir:      dir,
	}
	for _, avatar := range avatars {
		player.Avatars = append(player.Avatars, avatar.ToNetPlayerAvatar())
	}

	return player, nil
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

	// drop player current exp
	deathLossExp := lvSetting.DeathExpLoss
	if sceneData.Exp > deathLossExp {
		p.DeductExp(userId, deathLossExp)
		return nil
	}

	//  player max level item socket  lv -1
	maxLvSocket, err := p.playerMaxLevelItemSlot(userId)
	if err == nil && maxLvSocket.Level > 1 {
		p.setPlayerItemSlotLevel(userId, maxLvSocket.Position, maxLvSocket.Level-1, true)
	}

	return err
}
