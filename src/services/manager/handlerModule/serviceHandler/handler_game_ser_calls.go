package serviceHandler

import (
	"fmt"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
)

func GRPCServiceRegisterHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.ServiceRegisterOutput{ManagerAt: time_helper.NowUTCMill()}
	result := &module.ModuleEventResult{}
	defer func() {
		if result.GetError() != nil {
			serviceLog.Error(result.GetError().Error())
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.ServiceRegisterInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		result.SetError(err)
		return
	}

	service := controller.ToServiceData(input.Service)
	ctlModel, _ := controller.GetControllerModel()
	if _, exist := ctlModel.GetAliveServiceByType(
		service.ServiceType, service.SceneSerSubType, service.MapId, service.OwnerId,
	); !exist {
		// 玩家私有的服务 && 第一次注册时 发布启动完成事件
		ctlModel.GrpcCallPrivateSerStarted(&service)
	}

	serviceLog.Debug(
		"register service [%s],[%v],[%v],[%d] success",
		service.AppId, service.ServiceType, service.SceneSerSubType, service.OwnerId,
	)

	ctlModel.RegisterService(service)
	output.Success = true
	output.RegisterAt = input.RegisterAt
}

func GRPCServiceSelectHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.ManagerActionSelectServiceOutput{}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrorMessage != "" {
			output.ErrorCode = 30001
			serviceLog.Error(output.ErrorMessage)
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.ManagerActionSelectServiceInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrorMessage = err.Error()
		result.SetError(err)
		return
	}

	ctlModel, _ := controller.GetControllerModel()
	if serviceData, _ := ctlModel.GetAliveServiceByType(
		input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId,
	); serviceData == nil {
		output.ErrorMessage = fmt.Sprintf("Service [%v][%d]not found", input.ServiceType, input.MapId)
	} else {
		output.Service = serviceData.ToGrpcService()
	}

	serviceLog.Debug(
		"select service [%v],[%v],[%d],[%d] response service data: %+v",
		input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId,
		output.Service,
	)
}

func GRPCMultiSelectServiceHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.MultiSelectServiceOutput{}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrorMessage != "" {
			output.ErrorCode = 30003
			serviceLog.Error(output.ErrorMessage)
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MultiSelectServiceInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrorMessage = err.Error()
		result.SetError(err)
		return
	}

	ctlModel, _ := controller.GetControllerModel()
	allService := ctlModel.AllServices()
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

func GRPCServiceStartHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.StartServiceOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			serviceLog.Error(output.ErrMsg)
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.StartServiceInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.Success = false
		output.ErrMsg = err.Error()
		return
	}

	ctlModel, _ := controller.GetControllerModel()

	ser, exist := ctlModel.GetAliveServiceByType(
		input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId,
	)
	if exist { // 服务已启动
		ctlModel.GrpcCallPrivateSerStarted(ser)
	} else {
		if _, err = ctlModel.StartUserPrivateService(
			input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId,
		); err != nil {
			output.Success = false
			output.ErrMsg = err.Error()
		}
	}
}
