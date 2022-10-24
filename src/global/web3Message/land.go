package message

import (
	"fmt"
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
func ToWeb3LandData(l *proto.LandData) LandData {
	return LandData{
		Id:        int(l.Id),
		OccupyAt:  int(l.OccupyAt),
		OwnerId:   fmt.Sprint(l.Owner),
		TimeoutAt: int(l.TimeoutAt),
		X:         float64(l.X),
		Y:         float64(l.Y),
		Z:         float64(l.Z),
	}
}
