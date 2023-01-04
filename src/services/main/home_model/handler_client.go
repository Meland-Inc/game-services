package home_model

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func (p *HomeModel) clientMsgHandler(env *component.ModelEventReq, curMs int64) {
	bs, ok := env.Msg.([]byte)
	serviceLog.Info("client msg: %s, [%v]", bs, ok)
	if !ok {
		serviceLog.Error("client msg to string failed: %v", bs)
		return
	}

	serviceLog.Info("main service received clientPbMsg data: %v", string(bs))

	input := &methodData.PullClientMessageInput{}
	err := grpcNetTool.UnmarshalGrpcData(bs, input)
	if err != nil {
		serviceLog.Error("client msg input Unmarshal error: %v", err)
		return
	}

	agent := userAgent.GetOrStoreUserAgent(input)

	msg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	if err != nil {
		serviceLog.Error("Unmarshal Envelope fail err: %+v", err)
		return
	}

	switch proto.EnvelopeType(input.MsgId) {
	case proto.EnvelopeType_QueryGranary:
		p.QueryGranaryHandler(agent, input, msg)
	case proto.EnvelopeType_GranaryCollect:
		p.GranaryCollectHandler(agent, input, msg)

	}

}

func (p *HomeModel) QueryGranaryHandler(
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

	rows, err := p.GetGranaryRows(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, row := range rows {
		res.Items = append(res.Items, row.ToProtoData())
	}
}

func (p *HomeModel) GranaryCollectHandler(
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

	rows, err := p.GetGranaryRows(input.UserId)
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
	p.ClearGranaryRecord(input.UserId)
	p.BroadCastUpAllGranary(input.UserId)
}
