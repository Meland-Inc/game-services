package controller

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func (this *ControllerModel) serviceRecordByType(sType proto.ServiceType) (*ServiceRecord, bool) {
	iRecord, exist := this.controller.Load(sType)
	if !exist {
		return nil, false
	}
	record, ok := iRecord.(*ServiceRecord)
	if !ok {
		serviceLog.Error("interface to *ServiceRecord fail")
		return nil, false
	}
	return record, true
}

func (this *ControllerModel) RegisterService(service ServiceData) {
	record, ok := this.serviceRecordByType(service.ServiceType)
	if !ok {
		record = NewServiceRecord(service.ServiceType)
		this.controller.Store(service.ServiceType, record)
	}
	record.UpdateOrAddServiceRecord(service)
}

func (this *ControllerModel) DestroyService(service ServiceData) {
	record, ok := this.serviceRecordByType(service.ServiceType)
	if !ok {
		return
	}
	record.RemoveServiceRecord(service.AppId)
}

func (this *ControllerModel) GetAliveServiceByType(
	sType proto.ServiceType,
	sceneSubType proto.SceneServiceSubType,
	mapId int32,
	ownerId int64,
) (*ServiceData, bool) {
	record, ok := this.serviceRecordByType(sType)
	if !ok {
		return nil, false
	}
	return record.GetAliveService(mapId, sceneSubType, ownerId)
}

func (this *ControllerModel) AllServices() (services []ServiceData) {
	this.controller.Range(func(key, value interface{}) bool {
		if record, ok := value.(*ServiceRecord); ok {
			for _, s := range record.Services {
				services = append(services, s)
			}
		}
		return true
	})
	return services
}

func (this *ControllerModel) checkAndRemoveTimeOutSer(curMs int64) {
	this.controller.Range(func(key, value interface{}) bool {
		if record, ok := value.(*ServiceRecord); ok {
			record.checkAndRemoveTimeoutSer(curMs)
		}
		return true
	})
}

func (this *ControllerModel) PrintAllServices() {
	this.controller.Range(func(key, value interface{}) bool {
		if record, ok := value.(*ServiceRecord); ok {
			for appId, s := range record.Services {
				serviceLog.Info("serviceType[%v], appId[%d], data:%+v", key.(proto.ServiceType), appId, s)
			}
		}
		return true
	})
	serviceLog.Info("----------------------------------------------------------------")
}
