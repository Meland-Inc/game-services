package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func TaskRewardHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.TaskRewardResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40005 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("TaskReward res err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_TaskRewardResponse{TaskRewardResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetTaskRewardRequest()
	if req == nil {
		respMsg.ErrorMessage = "TaskReward request is nil"
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		respMsg.ErrorMessage = "TaskReward taskModel not found"
		return
	}

	task, err := taskModel.TaskReward(input.UserId, req.TaskListKind)
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
