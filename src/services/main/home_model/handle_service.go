package home_model

import (
	"game-message-core/grpc/methodData"
	"game-message-core/grpc/pubsubEventData"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/dapr/go-sdk/service/common"
)

func (p *HomeModel) GRPCGetHomeDataHandler(env *component.ModelEventReq, curMs int64) {
	inputBs, ok := env.Msg.([]byte)
	serviceLog.Debug("received GetHomeData : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("GetHomeData to string failed: %s", inputBs)
		return
	}

	output := &methodData.MainServiceActionGetHomeDataOutput{Success: true}
	result := &component.ModelEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
		}
		serviceLog.Debug("getHomeData output = %+v", output)
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionGetHomeDataInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	homeData, err := p.GetUserHomeData(input.UserId)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}
	output.UserId = input.UserId
	output.Data = ToGrpcHomeData(*homeData)
}

// ------------- pubsub event -------------

func (p *HomeModel) GRPCSaveHomeDataEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
	serviceLog.Info("SaveHomeDataEvent : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("SaveHomeDataEvent to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.SaveHomeEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("SaveHomeDataEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("Receive SaveHomeDataEvent: %+v", input)
	if err = p.UpdateUserHomeData(input.UserId, input.Data); err != nil {
		serviceLog.Error("SaveHomeDataEvent up user home data failed err: %v ", err)
	}
}

func (p *HomeModel) GRPCGranaryStockpileEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
	serviceLog.Info("GranaryStockpile Event : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("GranaryStockpile Event to TopicEvent failed: %v", msg)
		return
	}

	input := &pubsubEventData.GranaryStockpileEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("GranaryStockpileEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	serviceLog.Info("Receive GranaryStockpileEvent: %+v", input)
	upTm := time.UnixMilli(input.MsgVersion).UTC()
	for _, it := range input.Items {
		err := p.TryAddGranaryRecord(input.HomeOwner, it.Cid, it.Num, it.Quality, upTm, input.OccupantId, input.OccupantName)
		if err != nil {
			serviceLog.Error(err.Error())
		}
	}
}
