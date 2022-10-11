package chatModel

import "game-message-core/proto"

// 视野网格 直径
const VIEW_GRID_WITH float64 = 20.0 // 20M

type GridPos int64

const (
	GRID_POS_LEFT_UP GridPos = iota
	GRID_POS_UP
	GRID_POS_RIGHT_UP
	GRID_POS_LEFT
	GRID_POS_CENTRE
	GRID_POS_RIGHT
	GRID_POS_LEFT_DOWN
	GRID_POS_DOWN
	GRID_POS_RIGHT_DOWN
)

var NearPosOffsets = map[GridPos]proto.Vector3{
	GRID_POS_LEFT_UP:    proto.Vector3{X: -1, Y: 0, Z: 1},
	GRID_POS_UP:         proto.Vector3{X: 0, Y: 0, Z: 1},
	GRID_POS_RIGHT_UP:   proto.Vector3{X: 1, Y: 0, Z: 1},
	GRID_POS_RIGHT:      proto.Vector3{X: 1, Y: 0, Z: 0},
	GRID_POS_RIGHT_DOWN: proto.Vector3{X: 1, Y: 0, Z: -1},
	GRID_POS_DOWN:       proto.Vector3{X: 0, Y: 0, Z: -1},
	GRID_POS_LEFT_DOWN:  proto.Vector3{X: -1, Y: 0, Z: -1},
	GRID_POS_LEFT:       proto.Vector3{X: -1, Y: 0, Z: 0},
}

/// 根据grid pos 计算grid id, 注意gridPos 不是场景坐标
func CalGridId(gridPosX int64, gridPosZ int64) int64 {
	// 例： 100333000004 PosXY 各占6位 总共12位
	// 100333 第一位 1 代表负数  333 为X坐标 中间用0填充(1*100000+333)
	// 000004 第一位 0 代表正数  4   为Y坐标 中间用0填充(0*100000+333)
	var idTemplate int64 = 100000
	var xTag int64 = 0
	if gridPosX < 0 {
		xTag = 1
		gridPosX = -gridPosX
	}
	var zTag int64 = 0
	if gridPosZ < 0 {
		zTag = 1
		gridPosZ = -gridPosZ
	}
	xOffset := (xTag * idTemplate) + gridPosX
	zOffset := (zTag * idTemplate) + gridPosZ
	return (xOffset * idTemplate * 10) + zOffset
}
