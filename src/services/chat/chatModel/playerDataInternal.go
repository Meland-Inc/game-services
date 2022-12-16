package chatModel

import (
	"errors"
	"fmt"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func (p *ChatModel) createMapGrid(mapId int32) (*MapGrid, error) {
	if mapId == 0 {
		return nil, errors.New("mapId must be non-zero")
	}

	mapGrid := NewMapGrid(mapId)
	p.mapGrids[mapId] = mapGrid
	return mapGrid, nil
}

func (p *ChatModel) GetMapGrid(mapId int32) (*MapGrid, error) {
	if mapId == 0 {
		return nil, errors.New("mapId must be non-zero")
	}

	mapGrid, exist := p.mapGrids[mapId]
	if !exist {
		return nil, fmt.Errorf("map[%d] MapGrid not found", mapId)
	}
	return mapGrid, nil
}

// if not found player chat data return nil
func (p *ChatModel) GetPlayerChatData(userId int64) *PlayerChatData {
	data, _ := p.Players[userId]
	return data
}

func (p *ChatModel) RemovePlayerRecord(userId int64) {
	data, exist := p.Players[userId]
	if !exist {
		return
	}

	mapGrid, err := p.GetMapGrid(data.MapId)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	mapGrid.RemovePlayerGrid(userId)
	delete(p.Players, userId)
}

func (p *ChatModel) AddPlayerRecord(player *PlayerChatData) error {
	if player == nil {
		return errors.New("PlayerGridData must not be nil")
	}

	mapGrid, err := p.GetMapGrid(player.MapId)
	if mapGrid == nil {
		mapGrid, err = p.createMapGrid(player.MapId)
	}
	if err != nil {
		serviceLog.Error(err.Error())
		return err
	}

	mapGrid.addPlayerGrid(player)
	p.Players[player.UserId] = player
	return nil
}

func (p *ChatModel) UpdatePlayerRecord(newChatData *PlayerChatData) error {
	if newChatData == nil {
		return errors.New("newChatData must not be nil")
	}

	// // 新进入聊天
	data, exist := p.Players[newChatData.UserId]
	if !exist {
		return errors.New("newChatData not found")
	}

	// 切换地图
	if data.MapId != newChatData.MapId {
		preMapGrid, err := p.GetMapGrid(data.MapId)
		if err != nil {
			return err
		}
		nextMapGrid, err := p.GetMapGrid(newChatData.MapId)
		if err != nil {
			return err
		}
		preMapGrid.RemovePlayerGrid(newChatData.UserId)
		nextMapGrid.addPlayerGrid(newChatData)
		return nil
	}

	// 更新坐标 && 切换viewGrid
	mapGrid, err := p.GetMapGrid(newChatData.MapId)
	if err != nil {
		return err
	}
	mapGrid.UpdateAndAddPlayerGrid(newChatData)
	return nil
}

func (p *ChatModel) UpdateAndAddPlayerRecord(player *PlayerChatData) error {
	if player == nil {
		return errors.New("PlayerGridData must not be nil")
	}

	// 新进入聊天
	_, exist := p.Players[player.UserId]
	if !exist {
		return p.AddPlayerRecord(player)
	}
	return p.UpdatePlayerRecord(player)
}

func (p *ChatModel) OnPlayerEnterGame(env *pubsubEventData.UserEnterGameEvent) error {
	if env == nil {
		return errors.New(" enter game evn is nil")
	}

	playerChatData := NewPlayerChatData(
		env.UserId, env.Name, env.RoleIcon,
		env.MapId, env.X, env.Y, env.Z,
		env.SceneServiceAppId, env.AgentAppId, env.UserSocketId,
	)
	return p.UpdateAndAddPlayerRecord(playerChatData)
}

func (p *ChatModel) OnUpdatePlayerData(env *pubsubEventData.SavePlayerEventData) {
	playerChatData, exist := p.Players[env.UserId]
	if !exist {
		return
	}

	newData := NewPlayerChatData(
		playerChatData.UserId, playerChatData.Name, playerChatData.RoleIcon,
		env.FormService.MapId, env.PosX, env.PosY, env.PosZ,
		playerChatData.SceneServiceAppId,
		playerChatData.AgentAppId,
		playerChatData.UserSocketId,
	)
	p.UpdatePlayerRecord(newData)
}

func (p *ChatModel) OnPlayerLeaveGame(userId int64) error {
	if userId == 0 {
		return errors.New("leave game evn userId is zero")
	}
	p.RemovePlayerRecord(userId)
	return nil
}
