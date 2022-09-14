package controller

import (
	"game-message-core/proto"
	"sync"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

var controller *ServiceController

func GetInstance() *ServiceController {
	if controller == nil {
		NewServiceController()
	}
	return controller
}

func NewServiceController() *ServiceController {
	controller = &ServiceController{}
	return controller
}

type ServiceController struct {
	controller sync.Map
}

func (this *ServiceController) serviceRecordByType(sType proto.ServiceType) (*ServiceRecord, bool) {
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

func (this *ServiceController) RegisterService(service ServiceData) {
	if service.UpdatedAt == 0 {
		service.UpdatedAt = time_helper.NowUTCMill()
	}
	if service.CreatedAt == 0 {
		service.CreatedAt = time_helper.NowUTCMill()
	}

	record, ok := this.serviceRecordByType(service.ServiceType)
	if !ok {
		record = NewServiceRecord(service.ServiceType)
		this.controller.Store(service.ServiceType, record)
	}
	record.UpdateOrAddServiceRecord(service)
}

func (this *ServiceController) DestroyService(service ServiceData) {
	record, ok := this.serviceRecordByType(service.ServiceType)
	if !ok {
		return
	}
	record.RemoveServiceRecord(service.Id)
}

func (this *ServiceController) GetAliveServiceByType(sType proto.ServiceType) (*ServiceData, bool) {
	record, ok := this.serviceRecordByType(sType)
	if !ok {
		return nil, false
	}
	return record.GetAliveService()
}

func (this *ServiceController) AllServices() (services []ServiceData) {
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
