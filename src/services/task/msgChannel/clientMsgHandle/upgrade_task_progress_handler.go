package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func UpgradeTaskProgressHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.UpgradeTaskProgressResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40004 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("UpgradeTaskProgress res err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_UpgradeTaskProgressResponse{UpgradeTaskProgressResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetUpgradeTaskProgressRequest()
	if req == nil {
		respMsg.ErrorMessage = "UpgradeTaskProgress request is nil"
		return
	}

	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		respMsg.ErrorMessage = "UpgradeTaskProgress taskModel not found"
		return
	}

	if _, err := taskModel.UpGradeTaskProgress(
		input.UserId, req.TaskListKind, req.Items, req.Pos, req.Quiz, 0, 0, 0, 0,
	); err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}
