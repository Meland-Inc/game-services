package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/task/msgChannel/clientMsgHandle"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerClientMsgHandler() {
	ch.clientMsgHandler[proto.EnvelopeType_SelfTasks] = clientMsgHandle.SelfTasksHandler
	ch.clientMsgHandler[proto.EnvelopeType_AcceptTask] = clientMsgHandle.AcceptTaskHandler
	ch.clientMsgHandler[proto.EnvelopeType_AbandonmentTask] = clientMsgHandle.AbandonmentTaskHandler
	ch.clientMsgHandler[proto.EnvelopeType_UpgradeTaskProgress] = clientMsgHandle.UpgradeTaskProgressHandler
	ch.clientMsgHandler[proto.EnvelopeType_TaskReward] = clientMsgHandle.TaskRewardHandler
	ch.clientMsgHandler[proto.EnvelopeType_TaskListReward] = clientMsgHandle.TaskListRewardHandler
}
