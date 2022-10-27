package message

import (
	"game-message-core/proto"

	"github.com/spf13/cast"
)

func ToProtoLandData(l LandData) *proto.LandData {
	return &proto.LandData{
		Id:        int32(l.Id),
		OccupyAt:  int32(l.OccupyAt),
		Owner:     cast.ToInt64(l.OwnerId),
		TimeoutAt: int32(l.TimeoutAt),
		X:         float32(l.X),
		Y:         float32(l.Y),
		Z:         float32(l.Z),
	}
}
