package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	land_model "github.com/Meland-Inc/game-services/src/services/main/landModel"
)

func getMapLandRecordByUser(userId int64) (*land_model.MapLandDataRecord, error) {
	dataModel, err := land_model.GetLandModel()
	if err != nil {
		return nil, err
	}
	return dataModel.GetMapLandRecordByUser(userId)
}

func queryLandsGroupingResponse(
	input *methodData.PullClientMessageInput,
	agent *userAgent.UserAgentData,
	lands []*proto.LandData,
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
		ResponseClientMessage(agent, input, msg)
	}
}

func QueryLandsHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.QueryLandsResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22001 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_QueryLandsResponse{QueryLandsResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Query Lands Invalid User ID"
		return
	}

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	lands, err := mapLandRecord.AllLandData()
	serviceLog.Debug("QueryLandsHandler land length = %+v", len(lands))
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	queryLandsGroupingResponse(input, agent, lands)
}

func OccupyLandHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.OccupyLandResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22002 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_OccupyLandResponse{OccupyLandResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Occupy Land Invalid User ID"
		return
	}

	req := msg.GetOccupyLandRequest()
	if req == nil {
		serviceLog.Error("main service Occupy Lands request is nil")
		return
	}

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = mapLandRecord.OccupyLand(input.UserId, req.LandId, req.X, req.Z)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func BuildHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.BuildResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22003 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_BuildResponse{BuildResponse: res}
		ResponseClientMessage(agent, input, respMsg)
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

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
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

func RecyclingHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.RecyclingResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22004 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_RecyclingResponse{RecyclingResponse: res}
		ResponseClientMessage(agent, input, respMsg)
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

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
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

func ChargedHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.ChargedResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_ChargedResponse{ChargedResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "land Build Charged Invalid User ID"
		return
	}

	req := msg.GetChargedRequest()
	if req == nil {
		serviceLog.Error("main service Build Charged request is nil")
		return
	}

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
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

func HarvestHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.HarvestResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_HarvestResponse{HarvestResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "land Build Harvest Invalid User ID"
		return
	}

	req := msg.GetHarvestRequest()
	if req == nil {
		serviceLog.Error("main service Build Harvest request is nil")
		return
	}

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
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

func CollectionHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.CollectionResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22007 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_CollectionResponse{CollectionResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "land Build Collection Invalid User ID"
		return
	}

	req := msg.GetCollectionRequest()
	if req == nil {
		serviceLog.Error("main service Build Collection request is nil")
		return
	}

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
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

func SelfNftBuildsHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.SelfNftBuildsResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 22008 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_SelfNftBuildsResponse{SelfNftBuildsResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "SelfNftBuilds Invalid User ID"
		return
	}

	mapLandRecord, err := getMapLandRecordByUser(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	builds := mapLandRecord.GetUserNftBuilds(input.UserId)
	serviceLog.Debug("SelfNftBuildsHandler  buildIds length=[%d]", len(builds))
	for _, build := range builds {
		res.Builds = append(res.Builds, build.ToProtoData())
	}
}
