package chatModel

import (
	"errors"
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func (p *ChatModel) OnReceiveChatPbMsg(userId, msgId int64, chatMsg *proto.Envelope) error {
	if userId == 0 || msgId == 0 || chatMsg == nil {
		return errors.New("invalid chat message")
	}

	req := chatMsg.GetSendChatMessageRequest()
	if req == nil {
		return errors.New("chat msg request is nil")
	}

	context := p.checkChatContent(req.Content)
	if len(context) == 0 {
		return errors.New("chat content is empty")
	}

	chatData := &proto.ChatMessage{
		MsgId:        msgId,
		SenderId:     userId,
		ChatType:     req.ChatType,
		Content:      context,
		ReceiverUser: req.ReceiverPlayerId,
	}
	p.broadcastChatMsg(chatData)
	return nil
}

func (p *ChatModel) checkChatContent(content string) string {
	// content = m.filter.Replace(content, '*')
	return content
}

func (p *ChatModel) broadcastChatMsg(chatData *proto.ChatMessage) error {
	playerChatData := p.Players[chatData.SenderId]
	if playerChatData == nil {
		return fmt.Errorf("user [%v] chat data not found", chatData.SenderId)
	}
	curMs := time_helper.NowUTCMill()
	nextAt, _ := playerChatData.ChatCDs[chatData.ChatType]
	if nextAt > curMs {
		return fmt.Errorf("chat [%v] in CD")
	}

	chatData.SenderName = playerChatData.Name
	chatData.SenderIcon = playerChatData.RoleIcon

	switch chatData.ChatType {
	case proto.ChatChannelType_ChatChannelTypeSystem: // 系统
	case proto.ChatChannelType_ChatChannelTypeWorld: // 世界
		return p.worldChat(playerChatData, chatData, curMs)
	case proto.ChatChannelType_ChatChannelTypeNear: // 附近
		return p.nearChat(playerChatData, chatData, curMs)
	case proto.ChatChannelType_ChatChannelTypePrivate: // 私聊
		return p.privateChat(playerChatData, chatData, curMs)
	}
	return nil
}

func (p *ChatModel) makeChatBroadCastPbMsg(chatDatas []*proto.ChatMessage) *proto.Envelope {
	return &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastChatMessages,
		Payload: &proto.Envelope_BroadCastChatMessagesResponse{
			BroadCastChatMessagesResponse: &proto.BroadCastChatMessagesResponse{
				Messages: chatDatas,
			},
		},
	}
}

func (p *ChatModel) nearChat(sender *PlayerChatData, chatData *proto.ChatMessage, curMs int64) error {
	msg := p.makeChatBroadCastPbMsg([]*proto.ChatMessage{chatData})
	sender.InGrid.BroadcastNearMessage(msg, sender.UserId)
	sender.UpChatCD(chatData.ChatType)
	return nil
}

func (p *ChatModel) privateChat(sender *PlayerChatData, chatData *proto.ChatMessage, curMs int64) error {
	receiverInfo := p.GetPlayerChatData(chatData.ReceiverUser)
	if receiverInfo == nil {
		return fmt.Errorf("receiver [%d] chat data not found", chatData.ReceiverUser)
	}
	sender.UpChatCD(chatData.ChatType)
	msg := p.makeChatBroadCastPbMsg([]*proto.ChatMessage{chatData})
	err := userAgent.BroadCastToClient(
		receiverInfo.AgentAppId,
		serviceCnf.GetInstance().ServerName,
		receiverInfo.UserId,
		receiverInfo.UserSocketId,
		msg,
	)
	return err
}

func (p *ChatModel) worldChat(sender *PlayerChatData, chatData *proto.ChatMessage, curMs int64) error {
	sender.UpChatCD(chatData.ChatType)
	msg := p.makeChatBroadCastPbMsg([]*proto.ChatMessage{chatData})
	agentList := make(map[string][]int64)
	for _, player := range p.Players {
		agentId := player.AgentAppId
		if _, exist := agentList[agentId]; exist {
			agentList[agentId] = append(agentList[agentId], player.UserId)
		} else {
			agentList[agentId] = []int64{player.UserId}
		}
	}
	serviceAppId := serviceCnf.GetInstance().ServerName
	for agentId, userIds := range agentList {
		err := userAgent.MultipleBroadCastToClient(agentId, serviceAppId, userIds, msg)
		if err != nil {
			serviceLog.Error(err.Error())
		}
	}
	return nil
}
