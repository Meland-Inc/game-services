package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func ItemGetGroupingResponse(userId int64, pbItems []*proto.Item) {
	serviceName := serviceCnf.GetInstance().ServerName
	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(userId)
	if !exist {
		serviceLog.Warning("user [%d] agent data not found", userId)
		return
	}

	addRes := &proto.BroadCastItemAddResponse{}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastItemAdd,
		Payload: &proto.Envelope_BroadCastItemAddResponse{
			BroadCastItemAddResponse: addRes,
		},
	}

	n := 9
	itemLength := len(pbItems)
	left := itemLength / n
	for i := 0; i <= left; i++ {
		begin := i * n
		max := begin + n
		if max > itemLength {
			max = itemLength
		}
		addRes.Items = pbItems[begin:max]
		agent.SendToPlayer(serviceName, msg)
	}
}

func ItemGetHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.ItemGetResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20002 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_ItemGetResponse{ItemGetResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	serviceLog.Info("main service userId[%v] get items ", input.UserId)

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item Get Invalid User ID"
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	playerItems, err := dataModel.GetPlayerItems(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	pbItems := []*proto.Item{}
	for _, it := range playerItems.Items {
		pbItems = append(pbItems, it.ToNetItem())
	}
	serviceLog.Info("main service userId[%v] itemLength[%v]", input.UserId, len(pbItems))

	if len(pbItems) < 9 {
		res.Items = pbItems
	} else {
		ItemGetGroupingResponse(input.UserId, pbItems)
	}
}

func ItemUseHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.ItemUseResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_ItemUseResponse{ItemUseResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item use Invalid User ID"
		return
	}

	req := msg.GetItemUseRequest()
	if req == nil {
		serviceLog.Error("main service use item request is nil")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	dataModel.UseItem(input.UserId, req.ItemId)
}
