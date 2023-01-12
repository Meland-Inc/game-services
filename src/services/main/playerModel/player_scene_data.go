package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"gorm.io/gorm"
)

func (p *PlayerDataModel) getBirthData() (mapId int32, pos proto.Vector3) {
	// TODO: 此处数据需要 从配置中获取， 目前缺失
	return 10001, proto.Vector3{X: 270, Y: 15, Z: 136}
}

func (p *PlayerDataModel) initPlayerSceneData(tx *gorm.DB, userId int64) (*dbData.PlayerSceneData, error) {
	defaultMap, defaultPos := p.getBirthData()
	data := &dbData.PlayerSceneData{
		UserId:      userId,
		Hp:          200,
		Level:       1,
		Exp:         0,
		MapId:       defaultMap,
		X:           defaultPos.X,
		Y:           defaultPos.Y,
		Z:           defaultPos.Z,
		DirX:        0,
		DirY:        0,
		DirZ:        1,
		LastLoginAt: time_helper.NowUTC(),
	}
	lvCnf := configData.ConfigMgr().RoleLevelCnf(data.Level)
	if lvCnf != nil {
		data.Hp = lvCnf.HpLimit
	} else {
		return nil, fmt.Errorf("role level[%v]config not found", data.Level)
	}

	err := tx.Create(data).Error
	return data, err
}

func (p *PlayerDataModel) GetPlayerSceneData(userId int64) (row *dbData.PlayerSceneData, err error) {
	row = &dbData.PlayerSceneData{}
	gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
		err = tx.Where("user_id = ?", userId).First(row).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			row, err = p.initPlayerSceneData(tx, userId)
		}
		return err
	})
	return row, err
}

func (p *PlayerDataModel) UpPlayerSceneData(data *dbData.PlayerSceneData) error {
	if data == nil {
		return fmt.Errorf("data is nil")
	}

	return gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(data).Error
	})
}
