package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func TaskListRewardHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.TaskListRewardResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40006 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("TaskListReward res err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_TaskListRewardResponse{TaskListRewardResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetTaskListRewardRequest()
	if req == nil {
		respMsg.ErrorMessage = "TaskListReward request is nil"
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		respMsg.ErrorMessage = "TaskListReward taskModel not found"
		return
	}

	task, err := taskModel.TaskListReward(input.UserId, req.Kind)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	var pbTlData *proto.TaskList
	if task != nil {
		pbTlData = task.ToPbData()
	}
	res.TaskListInfo = pbTlData
}
