package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func AbandonmentTaskHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.AbandonmentTaskResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40003 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("AbandonmentTask res err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_AbandonmentTaskResponse{AbandonmentTaskResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetAbandonmentTaskRequest()
	if req == nil {
		respMsg.ErrorMessage = "AbandonmentTask request is nil"
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		respMsg.ErrorMessage = "AbandonmentTask taskModel not found"
		return
	}

	task, err := taskModel.AbandonmentTask(input.UserId, req.Kind)
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
