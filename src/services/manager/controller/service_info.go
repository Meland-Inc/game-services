package controller

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

const ServiceTimeoutMs int64 = 1000 * 10 // 10 seconds is services timeout

type ServiceData struct {
	Id          int64             `json:"id"`
	Name        string            `json:"name"`
	AppId       string            `json:"appId"`
	ServiceType proto.ServiceType `json:"serviceType"`
	Host        string            `json:"host"`
	Port        int32             `json:"port"`
	MapId       int32             `json:"mapId"`
	Online      int32             `json:"online"`
	MaxOnline   int32             `json:"maxOnline"`
	CreateAt    int64             `json:"createdAt"`
	UpdateAt    int64             `json:"updatedAt"`
}

type ServiceRecord struct {
	ServiceType  proto.ServiceType     `json:"serviceType"`
	Services     map[int64]ServiceData `json:"services"`
	statusRecord *ServiceStatusRecord  `json:"-"`
}

func NewServiceRecord(serviceType proto.ServiceType) *ServiceRecord {
	return &ServiceRecord{
		ServiceType:  serviceType,
		Services:     make(map[int64]ServiceData),
		statusRecord: NewServiceStatusRecord(),
	}
}

func (sr *ServiceRecord) checkAlive(s ServiceData) bool {
	nowMs := time_helper.NowUTCMill()
	return nowMs < s.UpdateAt+ServiceTimeoutMs
}

func (sr *ServiceRecord) RemoveServiceRecord(serviceId int64) {
	delete(sr.Services, serviceId)
	sr.statusRecord.RemoveServiceStatusRecord(serviceId)
}

func (sr *ServiceRecord) AddServiceRecord(service ServiceData) bool {
	if _, exist := sr.Services[service.Id]; exist {
		return false
	}
	sr.Services[service.Id] = service
	sr.statusRecord.AddServiceStatusRecord(service.Id, service.Online, service.MaxOnline)
	return true
}

func (sr *ServiceRecord) UpdateOrAddServiceRecord(service ServiceData) {
	sr.Services[service.Id] = service
	sr.statusRecord.AddServiceStatusRecord(service.Id, service.Online, service.MaxOnline)
}

func (sr *ServiceRecord) GetAliveService(mapId int32) (s *ServiceData, exist bool) {
	if len(sr.Services) == 0 {
		return nil, false
	}

	findF := func(serIds []int64) (*ServiceData, bool) {
		for _, sId := range serIds {
			service, ok := sr.Services[sId]
			if !ok {
				continue
			}
			if service.ServiceType == proto.ServiceType_ServiceTypeScene && service.MapId != mapId {
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
