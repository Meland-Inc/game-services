package chatModel

import (
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

func (p *ChatModel) onStart() {
	onlinePlayerIds, onLinePlayers := p.loadOnlinePlayer()
	if len(onlinePlayerIds) < 1 {
		return
	}

	baseDatas, err := p.loadPlayerBaseData(onlinePlayerIds)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	sceneDatas, err := p.loadPlayerSceneData(onlinePlayerIds)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	for _, player := range onLinePlayers {
		baseData, exist := baseDatas[player.UserId]
		if !exist {
			serviceLog.Warning("user[%d] base data not found", player.UserId)
			continue
		}
		sceneData, exist := sceneDatas[player.UserId]
		if !exist {
			serviceLog.Warning("user[%d] scene data not found", player.UserId)
			continue
		}
		playerChatData := NewPlayerChatData(
			player.UserId,
			baseData.Name, baseData.RoleIcon,
			sceneData.MapId, sceneData.X, sceneData.Y, sceneData.Z,
			player.InSceneService, player.AgentAppId, player.SocketId,
		)
		p.UpdateAndAddPlayerRecord(playerChatData)
	}
}

func (p *ChatModel) loadOnlinePlayer() (playerIds []int64, players []dbData.LoginData) {
	players = []dbData.LoginData{}
	err := gameDB.GetGameDB().Find(&players).Error
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	playerIds = make([]int64, len(players))
	for idx, onlinePlayer := range players {
		playerIds[idx] = onlinePlayer.UserId
	}
	return
}

func (p *ChatModel) loadPlayerBaseData(
	playerIds []int64,
) (map[int64]dbData.PlayerBaseData, error) {
	arr := []dbData.PlayerBaseData{}
	err := gameDB.GetGameDB().Where("user_id IN ?", playerIds).Find(&arr).Error
	if err != nil {
		serviceLog.Error(err.Error())
		return nil, err
	}

	baseDatas := make(map[int64]dbData.PlayerBaseData)
	for _, player := range arr {
		baseDatas[player.UserId] = player
	}
	return baseDatas, nil
}

func (p *ChatModel) loadPlayerSceneData(
	playerIds []int64,
) (map[int64]dbData.PlayerSceneData, error) {
	arr := []dbData.PlayerSceneData{}
	err := gameDB.GetGameDB().Where("user_id IN ?", playerIds).Find(&arr).Error
	if err != nil {
		serviceLog.Error(err.Error())
		return nil, err
	}

	sceneDatas := make(map[int64]dbData.PlayerSceneData)
	for _, player := range arr {
		sceneDatas[player.UserId] = player
	}
	return sceneDatas, nil
}
