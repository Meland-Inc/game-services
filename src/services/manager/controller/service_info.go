package controller

import (
	"game-message-core/jsonData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

const ServiceTimeoutMs int64 = 1000 * 10 // 10 seconds is services timeout

type ServiceRecord struct {
	ServiceType  proto.ServiceType              `json:"serviceType"`
	Services     map[int64]jsonData.ServiceData `json:"services"`
	statusRecord *ServiceStatusRecord           `json:"-"`
}

func NewServiceRecord(serviceType proto.ServiceType) *ServiceRecord {
	return &ServiceRecord{
		ServiceType:  serviceType,
		Services:     make(map[int64]jsonData.ServiceData),
		statusRecord: NewServiceStatusRecord(),
	}
}

func (sr *ServiceRecord) checkAlive(s jsonData.ServiceData) bool {
	nowMs := time_helper.NowUTCMill()
	return nowMs < s.UpdatedAt+ServiceTimeoutMs
}

func (sr *ServiceRecord) RemoveServiceRecord(serviceId int64) {
	delete(sr.Services, serviceId)
	sr.statusRecord.RemoveServiceStatusRecord(serviceId)
}

func (sr *ServiceRecord) AddServiceRecord(service jsonData.ServiceData) bool {
	if _, exist := sr.Services[service.Id]; exist {
		return false
	}
	sr.Services[service.Id] = service
	sr.statusRecord.AddServiceStatusRecord(service.Id, service.Online, service.MaxOnline)
	return true
}

func (sr *ServiceRecord) UpdateOrAddServiceRecord(service jsonData.ServiceData) {
	sr.Services[service.Id] = service
	sr.statusRecord.AddServiceStatusRecord(service.Id, service.Online, service.MaxOnline)
}

func (sr *ServiceRecord) GetAliveService() (s *jsonData.ServiceData, exist bool) {
	if len(sr.Services) == 0 {
		return nil, false
	}

	findF := func(serIds []int64) (*jsonData.ServiceData, bool) {
		for _, sId := range serIds {
			service, ok := sr.Services[sId]
			if !ok {
				continue
			}
			if !sr.checkAlive(service) {
				serviceLog.Info("remove time service [%v][%v]", service.AppId, service.ServiceType)
				sr.RemoveServiceRecord(sId)
				continue
			}
			return &service, true
		}
		return nil, false
	}

	tarSerIds := sr.statusRecord.GetServicesByStatus(ServiceStatusNormal)
	if s, exist = findF(tarSerIds); exist {
		return s, exist
	}
	tarSerIds = sr.statusRecord.GetServicesByStatus(ServiceStatusFree)
	if s, exist = findF(tarSerIds); exist {
		return s, exist
	}
	tarSerIds = sr.statusRecord.GetServicesByStatus(ServiceStatusBusy)
	if s, exist = findF(tarSerIds); exist {
		return s, exist
	}
	return nil, false
}
