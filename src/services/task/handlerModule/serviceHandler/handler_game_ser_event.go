package serviceHandler

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/task/taskModel"
)

func GRPCUserEnterGameEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserEnterGameEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserEnterGameEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("Receive UserEnterGameEvent: %+v", input)

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if !exist {
		agentModel.AddUserAgentRecord(
			input.UserId,
			input.AgentAppId,
			input.UserSocketId,
			input.SceneServiceAppId,
		)
	} else {
		agent.TryUpdate(input.UserId, input.AgentAppId, input.UserSocketId, input.SceneServiceAppId)
	}
}

func GRPCUserLeaveGameEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserLeaveGameEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserLeaveGameEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive LeaveGame: %+v", input)

	agentModel := userAgent.GetUserAgentModel()
	agentModel.RemoveUserAgentRecord(input.UserId)
}

func GRPCUserChangeServiceEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserChangeServiceEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("UserChangeServiceEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("Receive UserChangeServiceEvent: %+v", input)
	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if exist {
		agent.TryUpdate(agent.UserId, agent.AgentAppId, agent.SocketId, input.ToService.AppId)
	} else {
		agent, _ = agentModel.AddUserAgentRecord(
			input.UserId,
			input.UserAgentAppId,
			input.UserSocketId,
			input.ToService.AppId,
		)
	}
}

func GRPCKillMonsterEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.KillMonsterEventData{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("KillMonsterEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive KillMonsterEvent: %+v", input)

	taskModel, _ := taskModel.GetTaskModel()
	if err := taskModel.KillMonsterHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		&proto.TaskOptionKillMonster{MonCid: input.MonsterCid, Num: 1},
	); err != nil {
		serviceLog.Error("task killMon handler err:%+v", err)
	}

	if len(input.DropList) > 0 {
		pickItems := []*proto.TaskOptionItem{}
		for _, drop := range input.DropList {
			pickItems = append(pickItems, &proto.TaskOptionItem{ItemCid: drop.Cid, Num: drop.Num})
		}
		if err := taskModel.GetItemHandler(
			input.UserId, proto.TaskListType_TaskListTypeUnknown, pickItems,
		); err != nil {
			serviceLog.Error("task get item handler err:%+v", err)
		}
	}
}

func GRPCPlayerUseItemEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserUseNFTEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("SavePlayerDataEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	serviceLog.Info("receive UserUseNFTEvent: %+v", input)

	taskModel, _ := taskModel.GetTaskModel()
	if err := taskModel.UseItemHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		&proto.TaskOptionItem{ItemCid: input.Cid, Num: input.Num},
	); err != nil {
		serviceLog.Error("task use item handler err:%+v", err)
	}
}

func GRPCUserLevelUpgradeEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.UserLevelUpgradeEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("PlayerDeathEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive UserLevelUpgradeEvent: %+v", input)

	taskModel, _ := taskModel.GetTaskModel()
	if err := taskModel.UserLevelHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		input.Level,
	); err != nil {
		serviceLog.Error("task user level handler err:%+v", err)
	}

}

func GRPCSlotLevelUpgradeEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.SlotLevelUpgradeEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("SlotLevelUpgradeEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive SlotLevelUpgradeEvent: %+v", input)

	taskModel, _ := taskModel.GetTaskModel()
	if err := taskModel.TargetSlotLevelHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		&proto.TaskOptionTargetSlotLevel{SlotPos: input.SlotPos, Level: input.Level},
	); err != nil {
		serviceLog.Error("task slot level upgrade err:%+v", err)
	}

	if err := taskModel.SlotLevelCountHandler(
		input.UserId, proto.TaskListType_TaskListTypeUnknown,
	); err != nil {
		serviceLog.Error("task slot level count err:%+v", err)
	}
}

func GRPCTaskFinishEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.TaskFinishEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("TaskFinishEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive TaskFinishEvent: %+v", input)

	taskModel, _ := taskModel.GetTaskModel()
	err = taskModel.TaskFinishCountHandler(input.UserId, input.TaskListType)
	if err != nil {
		serviceLog.Error("task TaskFinishEvent err:%+v", err)
	}

	getItems := []*proto.TaskOptionItem{}
	for _, item := range input.RewardItem {
		getItems = append(getItems, &proto.TaskOptionItem{
			ItemCid: item.Cid,
			Num:     item.Num,
		})
	}
	if err := taskModel.GetItemHandler(
		input.UserId, proto.TaskListType_TaskListTypeUnknown, getItems,
	); err != nil {
		serviceLog.Error("task TaskFinishEvent get item err:%+v", err)
	}

}

func GRPCTaskListFinishEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.TaskListFinishEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("TaskListFinishEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service receive TaskListFinishEvent: %+v", input)

	taskModel, _ := taskModel.GetTaskModel()

	if err := taskModel.TaskListTypeCountHandler(
		input.UserId,
		proto.TaskListType_TaskListTypeUnknown,
		input.TaskListType,
	); err != nil {
		serviceLog.Error("task TaskListFinishEvent err:%+v", err)
	}

	getItems := []*proto.TaskOptionItem{}
	for _, item := range input.RewardItem {
		getItems = append(getItems, &proto.TaskOptionItem{
			ItemCid: item.Cid,
			Num:     item.Num,
		})
	}
	if err := taskModel.GetItemHandler(
		input.UserId, proto.TaskListType_TaskListTypeUnknown, getItems,
	); err != nil {
		serviceLog.Error("task TaskFinishEvent get item err:%+v", err)
	}
}
