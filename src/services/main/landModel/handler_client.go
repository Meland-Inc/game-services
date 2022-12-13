package land_model

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

func (p *LandModel) clientMsgHandler(env *component.ModelEventReq, curMs int64) {
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
	case proto.EnvelopeType_QueryLands:
		p.QueryLandsHandler(agent, input, msg)
	case proto.EnvelopeType_Build:
		p.BuildHandler(agent, input, msg)
	case proto.EnvelopeType_Recycling:
		p.RecyclingHandler(agent, input, msg)
	case proto.EnvelopeType_MintBattery:
		p.MintBatteryHandler(agent, input, msg)
	case proto.EnvelopeType_Charged:
		p.ChargedHandler(agent, input, msg)
	case proto.EnvelopeType_Harvest:
		p.HarvestHandler(agent, input, msg)
	case proto.EnvelopeType_Collection:
		p.CollectionHandler(agent, input, msg)
	case proto.EnvelopeType_SelfNftBuilds:
		p.SelfNftBuildsHandler(agent, input, msg)

	}

}

func (p *LandModel) queryLandsGroupingResponse(
	input *methodData.PullClientMessageInput, agent *userAgent.UserAgentData, lands []*proto.LandData,
) {
	if agent == nil {
		serviceLog.Warning("ItemGetGroupingResponse user [%d] agent data not found", input.UserId)
		return
	}

	addRes := &proto.BroadCastInitLandResponse{}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastInitLand,
		Payload: &proto.Envelope_BroadCastInitLandResponse{
			BroadCastInitLandResponse: addRes,
		},
	}

	n := 300
	landLength := len(lands)
	left := landLength / n
	if landLength%n > 0 {
		left++
	}
	for i := 0; i < left; i++ {
		beginIdx := i * n
		endIdx := beginIdx + n
		if endIdx > landLength {
			endIdx = landLength
		}
		addRes.Lands = lands[beginIdx:endIdx]
		userAgent.ResponseClientMessage(agent, input, msg)
	}
}

func (p *LandModel) QueryLandsHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.QueryLandsResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22001 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_QueryLandsResponse{QueryLandsResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Query Lands Invalid User ID"
		return
	}

	mapLandRecord, err := p.GetMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	lands, err := mapLandRecord.AllLandData()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	p.queryLandsGroupingResponse(input, agent, lands)
}

func (p *LandModel) BuildHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.BuildResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22003 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_BuildResponse{BuildResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "land Build Invalid User ID"
		return
	}

	req := msg.GetBuildRequest()
	if req == nil {
		serviceLog.Error("main service land Build request is nil")
		return
	}

	mapLandRecord, err := p.GetMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	build, err := mapLandRecord.Build(input.UserId, req.NftId, req.Position, req.LandIds)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	res.Build = build.ToProtoData()
}

func (p *LandModel) RecyclingHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {

	res := &proto.RecyclingResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22004 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_RecyclingResponse{RecyclingResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Recycling Build Invalid User ID"
		return
	}

	req := msg.GetRecyclingRequest()
	if req == nil {
		serviceLog.Error("main service Recycling Build request is nil")
		return
	}

	mapLandRecord, err := p.GetMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = mapLandRecord.Recycling(input.UserId, req.BuildId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func (p *LandModel) MintBatteryHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.MintBatteryResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22009 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_MintBatteryResponse{MintBatteryResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "MintBattery Invalid User ID"
		return
	}

	req := msg.GetMintBatteryRequest()
	if req == nil {
		respMsg.ErrorMessage = "MintBattery request is nil"
		return
	}
	err := grpcInvoke.MintBattery(input.UserId, req.MintNum, req.GiftNum)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
	}
}

func (p *LandModel) ChargedHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.ChargedResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_ChargedResponse{ChargedResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "land Build Charged Invalid User ID"
		return
	}

	req := msg.GetChargedRequest()
	if req == nil {
		respMsg.ErrorMessage = "Build Charged request is nil"
		return
	}

	mapLandRecord, err := p.GetMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = mapLandRecord.BuildCharged(input.UserId, req.NftId, req.BuildId, req.Num)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func (p *LandModel) HarvestHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.HarvestResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_HarvestResponse{HarvestResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "land Build Harvest Invalid User ID"
		return
	}

	req := msg.GetHarvestRequest()
	if req == nil {
		respMsg.ErrorMessage = "Build Harvest request is nil"
		return
	}

	mapLandRecord, err := p.GetMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = mapLandRecord.Harvest(input.UserId, req.NftId, req.BuildId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func (p *LandModel) CollectionHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.CollectionResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22007 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_CollectionResponse{CollectionResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "land Build Collection Invalid User ID"
		return
	}

	req := msg.GetCollectionRequest()
	if req == nil {
		respMsg.ErrorMessage = "Build Collection request is nil"
		return
	}

	mapLandRecord, err := p.GetMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = mapLandRecord.Collection(input.UserId, req.NftId, req.BuildId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func (p *LandModel) SelfNftBuildsHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.SelfNftBuildsResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22008 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_SelfNftBuildsResponse{SelfNftBuildsResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "SelfNftBuilds Invalid User ID"
		return
	}

	mapLandRecord, err := p.GetMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	builds := mapLandRecord.GetUserNftBuilds(input.UserId)
	for _, build := range builds {
		res.Builds = append(res.Builds, build.ToProtoData())
	}
}
