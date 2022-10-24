package land_model

import (
	base_data "game-message-core/grpc/baseData"
	"game-message-core/proto"

	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/spf13/cast"
)

type NftBuildData struct {
	GameData *dbData.NftBuild
	Web3Data *message.BuildData
}

func NewNftBuildData(gameData dbData.NftBuild, web3Data message.BuildData) *NftBuildData {
	return &NftBuildData{
		GameData: &gameData,
		Web3Data: &web3Data,
	}
}

func (p *NftBuildData) GetOwner() int64                 { return cast.ToInt64(p.Web3Data.UserId) }
func (p *NftBuildData) GetBuildId() int64               { return int64(p.Web3Data.BuildId) }
func (p *NftBuildData) GetNftId() string                { return p.Web3Data.NftId }
func (p *NftBuildData) GetGameData() *dbData.NftBuild   { return p.GameData }
func (p *NftBuildData) GetWeb3Data() *message.BuildData { return p.Web3Data }

func (p *NftBuildData) SetWeb3Data(data *message.BuildData) {
	if data == nil {
		return
	}
	p.Web3Data = data
}

func (p *NftBuildData) InLandIds() (landIds []int32) {
	for _, id := range p.Web3Data.LandIds {
		landIds = append(landIds, int32(id))
	}
	return
}

func (p *NftBuildData) ToProtoData() *proto.NftBuild {
	pos := &proto.Vector3{X: p.GameData.X, Y: p.GameData.Y, Z: p.GameData.Z}
	dir := &proto.Vector3{X: p.GameData.DirX, Y: p.GameData.DirY, Z: p.GameData.DirZ}
	pbBuild := &proto.NftBuild{
		Id:                  p.GetBuildId(),
		Cid:                 p.GameData.Cid,
		FromNft:             p.Web3Data.NftId,
		Owner:               p.GetOwner(),
		LandIds:             p.InLandIds(),
		Position:            pos,
		Dir:                 dir,
		ElectricEnd:         int32(p.Web3Data.ElectricEnd),
		ProduceBeginAt:      int32(p.Web3Data.ProduceBeginAt),
		HarvestItemCount:    int32(p.Web3Data.HarvestItemCount),
		CollectionItemCount: int32(p.Web3Data.CollectionItemCount),
		CollectionAt:        int32(p.Web3Data.CollectionAt),
	}
	return pbBuild
}

func (p *NftBuildData) ToGrpcData() base_data.GrpcNftBuild {
	pos := base_data.GrpcVector3{X: p.GameData.X, Y: p.GameData.Y, Z: p.GameData.Z}
	dir := base_data.GrpcVector3{X: p.GameData.DirX, Y: p.GameData.DirY, Z: p.GameData.DirZ}
	grpcBuild := base_data.GrpcNftBuild{
		Id:                  p.GetBuildId(),
		Cid:                 p.GameData.Cid,
		FromNft:             p.Web3Data.NftId,
		Owner:               p.GetOwner(),
		LandIds:             p.InLandIds(),
		Position:            pos,
		Dir:                 dir,
		ElectricEnd:         int32(p.Web3Data.ElectricEnd),
		ProduceBeginAt:      int32(p.Web3Data.ProduceBeginAt),
		HarvestItemCount:    int32(p.Web3Data.HarvestItemCount),
		CollectionItemCount: int32(p.Web3Data.CollectionItemCount),
		CollectionAt:        int32(p.Web3Data.CollectionAt),
	}
	return grpcBuild
}
