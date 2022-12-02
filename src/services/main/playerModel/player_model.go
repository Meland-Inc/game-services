package playerModel

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/shardCache"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

type PlayerDataModel struct {
	component.ModelBase

	cache    *shardCache.ShardedCache
	cacheTTL time.Duration
}

func NewPlayerModel() *PlayerDataModel {
	p := &PlayerDataModel{
		cacheTTL: time.Duration(10) * time.Minute,
		cache:    shardCache.NewSharded(shardCache.NoExpiration, time.Duration(60)*time.Second, 2^4),
	}
	p.InitBaseModel(p, component.MODEL_NAME_PLAYER_DATA)
	return p
}

func (p *PlayerDataModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("player model init service model manager is nil")
	}
	p.ModelBase.OnInit(modelMgr)
	return nil
}

func (p *PlayerDataModel) OnStart() error {
	return nil
}

func (p *PlayerDataModel) OnTick(utc time.Time) {
	p.ModelBase.OnTick(utc)
}

func (p *PlayerDataModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (p *PlayerDataModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (p *PlayerDataModel) OnEvent(env *component.ModelEventReq, curMs int64) {}

func (p *PlayerDataModel) Secondly(utc time.Time) {}

func (p *PlayerDataModel) Minutely(utc time.Time) {}

func (p *PlayerDataModel) Hourly(utc time.Time) {}

func (p *PlayerDataModel) Daily(utc time.Time) {}

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
