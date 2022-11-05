package clientMsgHandle

import (
	"errors"
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

	if err := tryUpTaskOptions(input.UserId, req); err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func tryUpTaskOptions(userId int64, req *proto.UpgradeTaskProgressRequest) error {
	taskModel, err := taskModel.GetTaskModel()
	if err != nil {
		return errors.New("UpgradeTaskProgress taskModel not found")
	}

	if req.Items != nil {
		return taskModel.HandInItemHandler(userId, req.TaskListKind, req.Items)
	}
	if req.Pos != nil {
		return taskModel.TargetPositionHandler(userId, req.TaskListKind, req.Pos)
	}
	if req.UseRecipe != nil {
		return taskModel.UseRecipeHandler(userId, req.TaskListKind, req.UseRecipe)
	}
	if req.CraftSkillLevel != nil {
		return taskModel.CraftSkillLevelHandler(userId, req.TaskListKind, req.CraftSkillLevel)
	}
	return nil
}
