package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/home_model"
)

func QueryGranaryHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.QueryGranaryResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 23001 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_QueryGranaryResponse{QueryGranaryResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Query granary Invalid User ID"
		return
	}

	granaryModel, err := home_model.GetHomeModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	rows, err := granaryModel.GetGranaryRows(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, row := range rows {
		res.Items = append(res.Items, row.ToProtoData())
	}
}

func GranaryCollectHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.GranaryCollectResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 23002 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_GranaryCollectResponse{GranaryCollectResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Query granary Invalid User ID"
		return
	}

	granaryModel, err := home_model.GetHomeModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	rows, err := granaryModel.GetGranaryRows(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, row := range rows {
		err := grpcInvoke.Web3MintNFT(input.UserId, row.ItemCid, row.Num, row.Quality, 0, 0)
		if err != nil {
			serviceLog.Error("Collect Item mint nft[%d] err: %v", row.ItemCid, err)
		} else {
			res.Items = append(res.Items, row.ToProtoData())
		}
	}
	granaryModel.ClearGranaryRecord(input.UserId)
	granaryModel.BroadCastUpAllGranary(input.UserId)
}
