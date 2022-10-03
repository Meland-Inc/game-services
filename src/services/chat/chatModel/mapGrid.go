package chatModel

import (
	"github.com/Meland-Inc/game-services/src/common/matrix"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

type MapGrid struct {
	MapId   int32
	Grids   map[int64]*ViewGrid
	Players map[int64]*PlayerChatData
}

func NewMapGrid(mapId int32) *MapGrid {
	return &MapGrid{
		MapId:   mapId,
		Grids:   make(map[int64]*ViewGrid),
		Players: make(map[int64]*PlayerChatData),
	}
}

/// 根据 grid pos 动态创建新的grid, 同时对near grid 添加引用关系  <returns>Created grid</returns>
func (mg *MapGrid) createGrid(gridPosX, gridPosZ, id int64) *ViewGrid {
	newGrid := NewViewGrid(mg.MapId, gridPosX, gridPosZ, id)
	newGrid.SetNearGrid(GRID_POS_CENTRE, newGrid)
	mg.Grids[newGrid.Id] = newGrid

	for nearPosType, offset := range NearPosOffsets {
		if offset.X == 0 && offset.Z == 0 {
			continue
		}

		nearGridX := gridPosX + int64(offset.X)
		nearGridZ := gridPosZ + int64(offset.Z)
		nearGridId := CalGridId(nearGridX, nearGridZ)
		nearGrid, exist := mg.Grids[nearGridId]
		if !exist {
			continue
		}

		newGrid.SetNearGrid(nearPosType, nearGrid)
		switch nearPosType {
		case GRID_POS_LEFT_UP:
			nearGrid.SetNearGrid(GRID_POS_RIGHT_DOWN, newGrid)

		case GRID_POS_UP:
			nearGrid.SetNearGrid(GRID_POS_DOWN, newGrid)

		case GRID_POS_RIGHT_UP:
			nearGrid.SetNearGrid(GRID_POS_LEFT_DOWN, newGrid)

		case GRID_POS_RIGHT:
			nearGrid.SetNearGrid(GRID_POS_LEFT, newGrid)

		case GRID_POS_RIGHT_DOWN:
			nearGrid.SetNearGrid(GRID_POS_LEFT_UP, newGrid)

		case GRID_POS_DOWN:
			nearGrid.SetNearGrid(GRID_POS_UP, newGrid)

		case GRID_POS_LEFT_DOWN:
			nearGrid.SetNearGrid(GRID_POS_RIGHT_UP, newGrid)

		case GRID_POS_LEFT:
			nearGrid.SetNearGrid(GRID_POS_RIGHT, newGrid)
		}
	}

	return newGrid
}

/// 根据场景中的3D坐标，查找其所在的grid
func (mg *MapGrid) gridByXZ(x, z float32) *ViewGrid {
	realX := matrix.Round(float64(x), 3)
	realZ := matrix.Round(float64(z), 3)

	gridPosX := int64(matrix.Floor(realX / VIEW_GRID_WITH))
	gridPosZ := int64(matrix.Floor(realZ / VIEW_GRID_WITH))
	gridId := CalGridId(gridPosX, gridPosZ)
	if grid, exist := mg.Grids[gridId]; exist {
		return grid
	}
	return mg.createGrid(gridPosX, gridPosZ, gridId)
}

/// 查询 场景坐标(XYZ) 所在的grid, 会动态创建grid
func (mg *MapGrid) GridByXYZ(x, y, z float32) *ViewGrid {
	return mg.gridByXZ(x, z)
}

/// 使用gridId 查询 grid, 查询失败 返回null
func (mg *MapGrid) GridById(gridId int64) *ViewGrid {
	grid, exist := mg.Grids[gridId]
	if !exist {
		return nil
	}
	return grid
}

func (mg *MapGrid) RemovePlayerGrid(userId int64) {
	if data, exist := mg.Players[userId]; exist {
		delete(mg.Players, userId)
		if data.InGrid != nil {
			data.InGrid.RemovePlayer(userId)
		}
	}
}

func (mg *MapGrid) UpdateAndAddPlayerGrid(player *PlayerChatData) {
	_, exist := mg.Players[player.UserId]
	if exist {
		mg.updatePlayerGrid(player)
	} else {
		mg.addPlayerGrid(player)
	}
}

func (mg *MapGrid) updatePlayerGrid(player *PlayerChatData) {
	data, exist := mg.Players[player.UserId]
	if !exist {
		return
	}
	if data.InGrid != nil {
		data.InGrid.RemovePlayer(player.UserId)
	}

	newGrid := mg.GridByXYZ(player.X, player.Y, player.Z)
	if newGrid == nil {
		serviceLog.Error(
			"mapId[%d] gridPos[X:%v, Y:%v, Z:%v] not found",
			mg.MapId, player.X, player.Y, player.Z,
		)
		return
	}

	newGrid.AddPlayer(data)
	data.InGrid = newGrid
	data.X = player.X
	data.Y = player.Y
	data.Z = player.Z
}

func (mg *MapGrid) addPlayerGrid(player *PlayerChatData) {
	newGrid := mg.GridByXYZ(player.X, player.Y, player.Z)
	if newGrid == nil {
		serviceLog.Error(
			"mapId[%d] gridPos[X:%v, Y:%v, Z:%v] not found",
			mg.MapId, player.X, player.Y, player.Z,
		)
		return
	}

	newGrid.AddPlayer(player)
	player.InGrid = newGrid
	mg.Players[player.UserId] = player
}
