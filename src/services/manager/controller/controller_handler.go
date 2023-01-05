package controller

import (
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/dapr/go-sdk/service/common"
)

func (p *ControllerModel) OnEvent(env contract.IModuleEventReq, curMs int64) {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.StackError("ControllerModel.onEvent err: %v", err)
		}
	}()

	switch env.GetEventType() {
	case string(grpc.ManagerServiceActionRegister):
		p.RegisterServiceHandler(env, curMs)
	case string(grpc.ManagerServiceActionSelectService):
		p.SelectServiceHandler(env, curMs)
	case string(grpc.ManagerServiceActionMultiSelectService):
		p.MultiSelectServiceHandler(env, curMs)
	case string(grpc.ManagerServiceActionStartService):
		p.StartServiceHandler(env, curMs)

	case string(grpc.SubscriptionEventServiceUnregister):
		p.UnregisterServiceEvent(env, curMs)
	}

}

func (p *ControllerModel) UnregisterServiceEvent(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().(*common.TopicEvent)
	if !ok {
		serviceLog.Error("service Unregister to string failed: %v", msg)
		return
	}

	input := &pubsubEventData.ServiceUnregisterEvent{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("ServiceUnregisterEvent Unmarshal fail err: %v", err)
		return
	}
	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service UnRegister: %v", input)
	p.DestroyService(ToServiceData(input.Service))
}

func (p *ControllerModel) RegisterServiceHandler(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().([]byte)
	// serviceLog.Info("service register : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("service register to string failed: %v", msg)
		return
	}

	output := &methodData.ServiceRegisterOutput{ManagerAt: time_helper.NowUTCMill()}
	result := &module.ModuleEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.ServiceRegisterInput{}
	err := grpcNetTool.UnmarshalGrpcData(msg, input)
	if err != nil {
		output.Success = false
		result.SetError(err)
		return
	}

	service := ToServiceData(input.Service)
	_, exist := p.GetAliveServiceByType(service.ServiceType, service.SceneSerSubType, service.MapId, service.OwnerId)
	if !exist {
		// 玩家私有的服务 && 第一次注册时 发布启动完成事件
		p.GrpcCallPrivateSerStarted(&service)
	}
	// serviceLog.Debug("register service success %+v", service)
	p.RegisterService(service)
	output.Success = true
	output.RegisterAt = input.RegisterAt
}

func (p *ControllerModel) SelectServiceHandler(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().([]byte)
	serviceLog.Info("select service: %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("select ser msg to string failed: %v", msg)
		return
	}

	output := &methodData.ManagerActionSelectServiceOutput{}
	result := &module.ModuleEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.ManagerActionSelectServiceInput{}
	err := grpcNetTool.UnmarshalGrpcData(msg, input)
	if err != nil {
		output.ErrorMessage = err.Error()
		output.ErrorCode = 30001
		result.SetError(err)
		return
	}

	if serviceData, _ := p.GetAliveServiceByType(
		input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId,
	); serviceData == nil {
		output.ErrorCode = 30002
		output.ErrorMessage = fmt.Sprintf("Service [%v][%d]not found", input.ServiceType, input.MapId)
	} else {
		output.Service = serviceData.ToGrpcService()
	}
}

func (p *ControllerModel) MultiSelectServiceHandler(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().([]byte)
	serviceLog.Info("multi service register: %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("multi select ser msg to string failed: %v", msg)
		return
	}

	output := &methodData.MultiSelectServiceOutput{}
	result := &module.ModuleEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MultiSelectServiceInput{}
	err := grpcNetTool.UnmarshalGrpcData(msg, input)
	if err != nil {
		output.ErrorMessage = err.Error()
		output.ErrorCode = 30003
		result.SetError(err)
		return
	}

	allService := p.AllServices()
	for _, s := range allService {
		if s.ServiceType != input.ServiceType {
			continue
		}
		if s.SceneSerSubType != input.SceneSerSubType {
			continue
		}
		if s.MapId != input.MapId {
			continue
		}
		output.Services = append(output.Services, s.ToGrpcService())
	}
}

func (p *ControllerModel) StartServiceHandler(env contract.IModuleEventReq, curMs int64) {
	msg, ok := env.GetMsg().([]byte)
	serviceLog.Info("received --start service : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("start service msg to string failed: %v", msg)
		return
	}

	output := &methodData.StartServiceOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.StartServiceInput{}
	err := grpcNetTool.UnmarshalGrpcData(msg, input)
	if err != nil {
		output.ErrMsg = err.Error()
		output.Success = false
		return
	}

	ser, exist := p.GetAliveServiceByType(input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId)
	serviceLog.Debug("start service exist[%v] ser = %+v", exist, ser)
	if exist { // 服务已启动
		p.GrpcCallPrivateSerStarted(ser)
	} else {
		if _, err = p.startUserPrivateService(
			input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId,
		); err != nil {
			output.Success = false
			output.ErrMsg = err.Error()
		}
	}
}
