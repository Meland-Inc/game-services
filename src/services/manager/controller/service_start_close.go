package controller

import (
	"fmt"
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
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
func (this *ControllerModel) startUserPrivateService(
	serType proto.ServiceType, subType proto.SceneServiceSubType,
	mapId int32, ownerId int64,
) (*ServiceData, error) {
	if mapId < 1 || ownerId < 1 {
		return nil, fmt.Errorf("invalid service mapId[%d] ownerId[%d]", mapId, ownerId)
	}

	/*
		   	SERVICE_SUB_TYPE=world   #(world | home | dungeon)
		   	HOME_OWNER=0             #( 0 |  home owner id)
		   	GAME_MAP_ID=10001
		   	APP_ID=game-service-world-${GAME_MAP_ID}-1
			APP_ID=game-service-dungeon-${GAME_MAP_ID}-N
			APP_ID=game-service-home-HomeOwnerId
	*/

	iSer, exist := this.startingPrivateSer.Load(ownerId)
	if exist {
		return iSer.(*ServiceData), nil
	}

	appId := fmt.Sprintf("game-service-home-%d", ownerId)
	if subType == proto.SceneServiceSubType_Dungeon {
		appId = fmt.Sprintf("game-service-dungeon-%d-%d", mapId, ownerId)
	}

	startSer := &ServiceData{
		AppId:           appId,
		ServiceType:     proto.ServiceType_ServiceTypeScene,
		SceneSerSubType: subType,
		OwnerId:         ownerId,
		CreateAt:        time_helper.NowUTCMill(),
	}
	this.AddStartingService(startSer)

	//// TODO ... CALL start service and wait start res
	// output, err := rpcCallStartSceneService(subType, ownerId, mapId, appId)
	// if err != nil {
	// 	return nil, err
	// }
	// if !output.Success {
	// 	return nil, fmt.Errorf(output.FailedReason)
	// }

	return startSer, nil

	// TODO ... 监听目标服务启动完成事件
	// 调用 this.RemoveStartingService(ownerId)
}

// 关闭私有(家园|副本)
func closeUserPrivateService(ser ServiceData) error {
	serviceLog.Info("close user private ser %+v", ser.AppId, ser.SceneSerSubType)
	//// TODO ... CALL close service and wait start res
	// output, err := rpcCallCloseSceneService(ser.AppId)
	// if err != nil {
	// 	return nil, err
	// }
	// if !output.Success {
	// 	return nil, fmt.Errorf(output.FailedReason)
	// }

	return nil
}
