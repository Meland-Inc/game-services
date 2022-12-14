package playerModel

import (
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"gorm.io/gorm"
)

func (p *PlayerDataModel) getBirthData() (mapId int32, pos proto.Vector3) {
	// TODO: 此处数据需要 从配置中获取， 目前缺失
	return 10001, proto.Vector3{X: 271, Y: 6, Z: 45}
}

func (p *PlayerDataModel) initPlayerSceneData(userId int64) (*dbData.PlayerSceneData, error) {
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
		BirthMapId:  defaultMap,
		BirthX:      defaultPos.X,
		BirthY:      defaultPos.Y,
		BirthZ:      defaultPos.Z,
		LastLoginAt: time_helper.NowUTC(),
	}
	lvCnf := configData.ConfigMgr().RoleLevelCnf(data.Level)
	if lvCnf != nil {
		data.Hp = lvCnf.HpLimit
	} else {
		serviceLog.Error("role level[%v]config not found", data.Level)
	}

	err := gameDB.GetGameDB().Create(data).Error
	return data, err
}

func (p *PlayerDataModel) GetPlayerSceneData(userId int64) (*dbData.PlayerSceneData, error) {
	cacheKey := p.getPlayerSceneDataKey(userId)
	iData, err := p.cache.GetOrStore(
		cacheKey,
		func() (interface{}, error) {
			data := &dbData.PlayerSceneData{}
			err := gameDB.GetGameDB().Where("user_id = ?", userId).First(data).Error
			if err != nil && err == gorm.ErrRecordNotFound {
				data, err = p.initPlayerSceneData(userId)
			}
			return data, err
		},
		p.cacheTTL,
	)
	if err != nil {
		return nil, err
	}

	p.cache.Touch(cacheKey, p.cacheTTL)
	return iData.(*dbData.PlayerSceneData), nil
}

func (p *PlayerDataModel) UpPlayerSceneData(
	userId int64,
	hp, level, exp, mapId int32,
	x, y, z, dirX, dirY, dirZ float32,
) error {
	if userId == 0 {
		return fmt.Errorf("invalid player scene data")
	}

	data, err := p.GetPlayerSceneData(userId)
	if err != nil {
		return err
	}

	err = gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
		data.Hp = hp
		data.Level = level
		data.Exp = exp
		data.MapId = mapId
		data.X = x
		data.Y = y
		data.Z = z
		data.DirX = dirX
		data.DirY = dirY
		data.DirZ = dirZ
		return tx.Save(data).Error
	})
	return err
}
