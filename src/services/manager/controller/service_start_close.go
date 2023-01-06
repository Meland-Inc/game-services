package controller

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (this *ControllerModel) GrpcCallPrivateSerStarted(ser *ServiceData) {
	if !IsUserPrivateSer(*ser) {
		return
	}

	// 保证在此时间内服务不会因为过期而关闭
	ser.UpdateUserLastOnlineAt()

	go func() {
		// 延时100MS 通知服务启动完成 以保证 grpc output消息先到达
		time.Sleep(time.Millisecond * 100)
		grpcPubsubEvent.RPCPubsubEventServiceStarted(ser.ToGrpcService())
	}()
}

func (this *ControllerModel) AddStartingService(ser *ServiceData) {
	this.startingPrivateSer.Store(ser.OwnerId, ser)
}

func (this *ControllerModel) RemoveStartingService(serOwner int64) {
	this.startingPrivateSer.Delete(serOwner)
}

// 因为启动需要等待消息回复，外部调用时最好使用 异步调用
func (this *ControllerModel) StartUserPrivateService(
	serType proto.ServiceType, subType proto.SceneServiceSubType, mapId int32, ownerId int64,
) (*ServiceData, error) {
	if mapId < 1 || ownerId < 1 {
		return nil, fmt.Errorf("invalid service mapId[%d] ownerId[%d]", mapId, ownerId)
	}

	iSer, exist := this.startingPrivateSer.Load(ownerId)
	if exist {
		return iSer.(*ServiceData), nil
	}

	appId, err := grpcInvoke.GRPCDynamicStartSceneService(subType, ownerId, mapId, 3000)
	if err != nil {
		return nil, err
	}

	startSer := &ServiceData{
		AppId:           appId,
		ServiceType:     proto.ServiceType_ServiceTypeScene,
		SceneSerSubType: subType,
		OwnerId:         ownerId,
		CreateAt:        time_helper.NowUTCMill(),
	}
	this.AddStartingService(startSer)
	return startSer, nil
	// 监听目标服务启动完成事件 调用 this.RemoveStartingService(ownerId)
}

// 关闭私有(家园|副本)
func closeUserPrivateService(ser ServiceData) error {
	if IsUserPrivateSer(ser) {
		return fmt.Errorf("not user private service, can not close")
	}

	serviceLog.Info("close user private ser %+v", ser.AppId, ser.SceneSerSubType)
	err := grpcPubsubEvent.Web3RPCEventCloseDynamicSceneService(ser.AppId)
	if err != nil {
		serviceLog.Error(err.Error())
	}
	return err
}
