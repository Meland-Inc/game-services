package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func SelfTasksHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.SelfTasksResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40001 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("SelfTasks res err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_SelfTasksResponse{SelfTasksResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetSelfTasksRequest()
	if req == nil {
		respMsg.ErrorMessage = "SelfTasks request is nil"
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		respMsg.ErrorMessage = "SelfTasks taskModel not found"
		return
	}

	tasks, err := taskModel.GetPlayerTask(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	res.Tasks = tasks.ToProtoData()
}
