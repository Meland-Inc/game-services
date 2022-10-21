package message

import (
	"game-message-core/proto"
)

func ToProtoLandData(l LandData) *proto.LandData {
	return &proto.LandData{
		Id:        int32(l.Id),
		OccupyAt:  int32(l.OccupyAt),
		Owner:     int64(l.Owner),
		TimeoutAt: int32(l.TimeoutAt),
		X:         int32(l.X),
		Z:         int32(l.Z),
	}
}
func ToWeb3LandData(l *proto.LandData) LandData {
	return LandData{
		Id:        int(l.Id),
		OccupyAt:  int(l.OccupyAt),
		Owner:     int(l.Owner),
		TimeoutAt: int(l.TimeoutAt),
		X:         int(l.X),
		Z:         int(l.Z),
	}
}
