package home_model

import (
	"game-message-core/grpc/methodData"
	"game-message-core/grpc/pubsubEventData"

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

	output := &methodData.MainServiceActionGetHomeDataOutput{}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionGetHomeDataInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		result.Err = err
		return
	}

	homeData, err := p.GetUserHomeData(input.UserId)
	if err != nil {
		output.Success = false
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
