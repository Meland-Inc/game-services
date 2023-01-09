package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func SelfTasksHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.SelfTasksResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40001 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("SelfTasks err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_SelfTasksResponse{SelfTasksResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetSelfTasksRequest()
	if req == nil {
		respMsg.ErrorMessage = "SelfTasks request is nil"
		return
	}

	taskModel, _ := taskModel.GetTaskModel()
	tasks, err := taskModel.GetPlayerTask(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	res.Tasks = tasks.ToProtoData()
}

func AcceptTaskHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.AcceptTaskResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40002 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("AcceptTask err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_AcceptTaskResponse{AcceptTaskResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetAcceptTaskRequest()
	if req == nil {
		respMsg.ErrorMessage = "AcceptTask request is nil"
		return
	}

	taskModel, _ := taskModel.GetTaskModel()
	taskListData, err := taskModel.AcceptTask(input.UserId, req.Kind)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	if taskListData != nil {
		res.TaskListInfo = taskListData.ToPbData()
	}
}

func AbandonmentTaskHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.AbandonmentTaskResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40003 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("AbandonmentTask err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_AbandonmentTaskResponse{AbandonmentTaskResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetAbandonmentTaskRequest()
	if req == nil {
		respMsg.ErrorMessage = "AbandonmentTask request is nil"
		return
	}

	taskModel, _ := taskModel.GetTaskModel()
	task, err := taskModel.AbandonmentTask(input.UserId, req.Kind)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	if task != nil {
		res.TaskListInfo = task.ToPbData()
	}
}

func UpgradeTaskProgressHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UpgradeTaskProgressResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40004 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("UpgradeTaskProgress err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_UpgradeTaskProgressResponse{UpgradeTaskProgressResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetUpgradeTaskProgressRequest()
	if req == nil {
		respMsg.ErrorMessage = "UpgradeTaskProgress request is nil"
		return
	}

	taskModel, _ := taskModel.GetTaskModel()
	if req.Items != nil {
		err := taskModel.HandInItemHandler(input.UserId, req.TaskListKind, req.Items)
		if err != nil {
			respMsg.ErrorMessage = err.Error()
			return
		}
	}
	if req.Pos != nil {
		err := taskModel.TargetPositionHandler(input.UserId, req.TaskListKind, req.Pos)
		if err != nil {
			respMsg.ErrorMessage = err.Error()
			return
		}
	}
	if req.UseRecipe != nil {
		err := taskModel.UseRecipeHandler(input.UserId, req.TaskListKind, req.UseRecipe)
		if err != nil {
			respMsg.ErrorMessage = err.Error()
			return
		}
	}
	if req.CraftSkillLevel != nil {
		err := taskModel.CraftSkillLevelHandler(input.UserId, req.TaskListKind, req.CraftSkillLevel)
		if err != nil {
			respMsg.ErrorMessage = err.Error()
			return
		}
	}
}

func TaskRewardHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.TaskRewardResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40005 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("TaskReward err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_TaskRewardResponse{TaskRewardResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetTaskRewardRequest()
	if req == nil {
		respMsg.ErrorMessage = "TaskReward request is nil"
		return
	}

	taskModel, _ := taskModel.GetTaskModel()
	taskListData, err := taskModel.TaskReward(input.UserId, req.TaskListKind)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	if taskListData != nil {
		res.TaskListInfo = taskListData.ToPbData()
	}
}

func TaskListRewardHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.TaskListRewardResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 40006 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("TaskReward err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_TaskListRewardResponse{TaskListRewardResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetTaskListRewardRequest()
	if req == nil {
		respMsg.ErrorMessage = "TaskListReward request is nil"
		return
	}

	taskModel, _ := taskModel.GetTaskModel()
	taskListData, err := taskModel.TaskListReward(input.UserId, req.Kind)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	if taskListData != nil {
		res.TaskListInfo = taskListData.ToPbData()
	}
}
