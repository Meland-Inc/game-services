package land_model

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/dapr/go-sdk/service/common"
)

func (p *LandModel) Web3MultiLandDataUpdateEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("Web3MultiLandDataUpdateEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("Web3MultiLandDataUpdateEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &message.MultiLandDataUpdateEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("Web3MultiLandDataUpdateEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiLandDataUpdateEvent: %+v", input)

	landGroup := make(map[int32][]*proto.LandData)
	for _, land := range input.Lands {
		mapId := int32(land.MapId)
		pbLandData := message.ToProtoLandData(land)
		if pbLandData == nil {
			continue
		}
		if _, exist := landGroup[mapId]; exist {
			landGroup[mapId] = append(landGroup[mapId], pbLandData)
		} else {
			landGroup[mapId] = []*proto.LandData{pbLandData}
		}
	}

	for mapId, upLands := range landGroup {
		mapRecord, err := p.GetMapLandRecord(mapId)
		if err != nil {
			serviceLog.Error("MultiLandDataUpdateEvent error: %v", err)
			continue
		}
		mapRecord.MultiUpdateLandData(upLands)
	}
}

func (p *LandModel) Web3MultiRecyclingEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("Web3MultiRecyclingEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("Web3MultiRecyclingEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &message.MultiRecyclingEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("Web3MultiRecyclingEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiRecyclingEvent: %+v", input)

	for _, info := range input.RecyclingInfos {
		mapRecord, err := p.GetMapLandRecord(int32(info.MapId))
		if err != nil {
			serviceLog.Error("MultiRecyclingEvent error: %v", err)
			return
		}

		err = mapRecord.OnReceiveRecyclingEvent(int64(info.BuildId))
		if err != nil {
			serviceLog.Error("MultiRecyclingEvent error: %v", err)
		}
	}
}

func (p *LandModel) Web3MultiBuildUpdateEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	serviceLog.Info("Web3MultiBuildUpdateEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("Web3MultiBuildUpdateEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &message.MultiBuildUpdateEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("Web3MultiBuildUpdateEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiBuildUpdateEvent: %+v", input)

	for _, build := range input.BuildDatas {
		mapRecord, err := p.GetMapLandRecord(int32(build.MapId))
		if err != nil {
			serviceLog.Error("MultiBuildUpdateEvent error: %v", err)
			continue
		}

		err = mapRecord.UpdateNftBuildWeb3Data(build)
		if err != nil {
			serviceLog.Error("MultiBuildUpdateEvent error: %v", err)
		}
	}
}
