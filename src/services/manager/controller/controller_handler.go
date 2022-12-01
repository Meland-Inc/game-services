package controller

import (
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/dapr/go-sdk/service/common"
)

func (p *ControllerModel) onEvent(env *component.ModelEvent, curMs int64) {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.StackError("ControllerModel.onEvent err: %v", err)
		}
	}()

	switch env.EventType {
	case string(grpc.ManagerServiceActionRegister):
		p.RegisterServiceHandler(env, curMs)
	case string(grpc.ManagerServiceActionSelectService):
		p.SelectServiceHandler(env, curMs)
	case string(grpc.ManagerServiceActionMultiSelectService):
		p.MultiSelectServiceHandler(env, curMs)

	case string(grpc.SubscriptionEventServiceUnregister):
		p.UnregisterServiceEvent(env, curMs)
	}

}

func (p *ControllerModel) UnregisterServiceEvent(env *component.ModelEvent, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
	serviceLog.Info("service Unregister : %s, [%v]", msg, ok)
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
	service := ServiceData{
		AppId:           input.Service.AppId,
		ServiceType:     input.Service.ServiceType,
		SceneSerSubType: input.Service.SceneSerSubType,
		OwnerId:         input.Service.Owner,
		Host:            input.Service.Host,
		Port:            input.Service.Port,
		MapId:           input.Service.MapId,
		Online:          input.Service.Online,
		MaxOnline:       input.Service.MaxOnline,
		CreateAt:        input.Service.CreatedAt,
		UpdateAt:        input.Service.UpdatedAt,
	}
	p.DestroyService(service)
}

func (p *ControllerModel) RegisterServiceHandler(env *component.ModelEvent, curMs int64) {
	msg, ok := env.Msg.([]byte)
	serviceLog.Info("service register : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("service register to string failed: %v", msg)
		return
	}

	output := &methodData.ServiceRegisterOutput{ManagerAt: time_helper.NowUTCMill()}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.ServiceRegisterInput{}
	err := grpcNetTool.UnmarshalGrpcData(msg, input)
	if err != nil {
		output.Success = false
		result.Err = err
		return
	}

	service := ServiceData{
		AppId:           input.Service.AppId,
		ServiceType:     input.Service.ServiceType,
		SceneSerSubType: input.Service.SceneSerSubType,
		OwnerId:         input.Service.Owner,
		Host:            input.Service.Host,
		Port:            input.Service.Port,
		MapId:           input.Service.MapId,
		Online:          input.Service.Online,
		MaxOnline:       input.Service.MaxOnline,
		CreateAt:        input.Service.CreatedAt,
		UpdateAt:        input.Service.UpdatedAt,
	}
	// serviceLog.Debug("register service success %+v", service)
	p.RegisterService(service)
	output.Success = true
	output.RegisterAt = input.RegisterAt

}

func (p *ControllerModel) SelectServiceHandler(env *component.ModelEvent, curMs int64) {
	msg, ok := env.Msg.([]byte)
	serviceLog.Info("select service: %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("select ser msg to string failed: %v", msg)
		return
	}

	output := &methodData.ManagerActionSelectServiceOutput{}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.ManagerActionSelectServiceInput{}
	err := grpcNetTool.UnmarshalGrpcData(msg, input)
	if err != nil {
		output.ErrorMessage = err.Error()
		output.ErrorCode = 30001
		result.Err = err
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

func (p *ControllerModel) MultiSelectServiceHandler(env *component.ModelEvent, curMs int64) {
	msg, ok := env.Msg.([]byte)
	serviceLog.Info("multi service register: %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("multi select ser msg to string failed: %v", msg)
		return
	}

	output := &methodData.MultiSelectServiceOutput{}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MultiSelectServiceInput{}
	err := grpcNetTool.UnmarshalGrpcData(msg, input)
	if err != nil {
		output.ErrorMessage = err.Error()
		output.ErrorCode = 30003
		result.Err = err
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
