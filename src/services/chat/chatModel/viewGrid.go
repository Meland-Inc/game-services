package chatModel

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

type ViewGrid struct {
	Id        int64
	PosX      int64
	PosZ      int64
	InMap     int32
	nearGrids map[GridPos]*ViewGrid
	players   map[int64]*PlayerChatData
}

func NewViewGrid(inMapId int32, gridPosX, gridPosZ, gridId int64) *ViewGrid {
	return &ViewGrid{
		InMap:     inMapId,
		Id:        gridId,
		PosX:      gridPosX,
		PosZ:      gridPosZ,
		nearGrids: make(map[GridPos]*ViewGrid),
		players:   make(map[int64]*PlayerChatData),
	}
}

func (g *ViewGrid) SetNearGrid(pos GridPos, nearG *ViewGrid) bool {
	if nearG == nil {
		return false
	}
	g.nearGrids[pos] = nearG
	return true
}

/// 查询以该grid为中心的 9宫格 grid
func (g *ViewGrid) NineGridsAndPos() map[GridPos]*ViewGrid {
	return g.nearGrids
}

func (g *ViewGrid) RangeNear(f func(grid *ViewGrid) bool) {
	for _, gv := range g.nearGrids {
		if !f(gv) {
			break
		}
	}
}

func (g *ViewGrid) RangePlayers(f func(*PlayerChatData) bool) {
	for _, p := range g.players {
		if !f(p) {
			break
		}
	}
}

func (g *ViewGrid) RangeNearPlayers(f func(*PlayerChatData) bool) {
	g.RangeNear(func(grid *ViewGrid) bool {
		if grid != nil {
			grid.RangePlayers(f)
		}
		return true
	})
}

func (g *ViewGrid) PlayerById(userId int64) (*PlayerChatData, bool) {
	baseData, exist := g.players[userId]
	return baseData, exist
}

func (g *ViewGrid) RemovePlayer(userId int64) {
	delete(g.players, userId)
}
func (g *ViewGrid) AddPlayer(playerData *PlayerChatData) {
	if playerData == nil {
		return
	}
	g.players[playerData.UserId] = playerData
}

/// <summary>
/// 对在grid中的玩家广播消息
/// </summary>
func (g *ViewGrid) Broadcast(msg *proto.Envelope, exceptEntity int64) {
	userIds := []int64{}
	g.RangeNearPlayers(func(player *PlayerChatData) bool {
		if player == nil || player.UserId == exceptEntity {
			return true
		}
		userIds = append(userIds, player.UserId)
		return true
	})

	err := userAgent.MultipleBroadCastToClient(serviceCnf.GetInstance().AppId, userIds, msg)
	if err != nil {
		serviceLog.Error(err.Error())
	}
}

/// 对以该grid作为中心点的九宫格中的所以grid广播消息
func (g *ViewGrid) BroadcastNearMessage(msg *proto.Envelope, exceptEntity int64) {
	g.RangeNear(func(grid *ViewGrid) bool {
		grid.Broadcast(msg, exceptEntity)
		return true
	})
}
