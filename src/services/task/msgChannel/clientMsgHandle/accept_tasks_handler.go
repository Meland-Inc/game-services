package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func AcceptTaskHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.AcceptTaskResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40002 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("AcceptTask res err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_AcceptTaskResponse{AcceptTaskResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetAcceptTaskRequest()
	if req == nil {
		respMsg.ErrorMessage = "AcceptTask request is nil"
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		respMsg.ErrorMessage = "AcceptTask taskModel not found"
		return
	}

	task, err := taskModel.AcceptTask(input.UserId, req.Kind)
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
