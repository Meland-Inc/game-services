package controller

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

const SERVICE_TIME_OUT_MS int64 = 1000 * 5 // 10 seconds is services timeout

type ServiceData struct {
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
	ServiceType  proto.ServiceType      `json:"serviceType"`
	Services     map[string]ServiceData `json:"services"`
	statusRecord *ServiceStatusRecord   `json:"-"`
}

func NewServiceRecord(serviceType proto.ServiceType) *ServiceRecord {
	return &ServiceRecord{
		ServiceType:  serviceType,
		Services:     make(map[string]ServiceData),
		statusRecord: NewServiceStatusRecord(),
	}
}

func (sr *ServiceRecord) checkAlive(s ServiceData) bool {
	nowMs := time_helper.NowUTCMill()
	return nowMs < s.UpdateAt+SERVICE_TIME_OUT_MS
}

func (sr *ServiceRecord) RemoveServiceRecord(appId string) {
	delete(sr.Services, appId)
	sr.statusRecord.RemoveServiceStatusRecord(appId)
}

func (sr *ServiceRecord) AddServiceRecord(service ServiceData) bool {
	if _, exist := sr.Services[service.AppId]; exist {
		return false
	}
	sr.Services[service.AppId] = service
	sr.statusRecord.AddServiceStatusRecord(service.AppId, service.Online, service.MaxOnline)
	return true
}

func (sr *ServiceRecord) UpdateOrAddServiceRecord(service ServiceData) {
	sr.Services[service.AppId] = service
	sr.statusRecord.AddServiceStatusRecord(service.AppId, service.Online, service.MaxOnline)
}

func (sr *ServiceRecord) GetAliveService(mapId int32) (s *ServiceData, exist bool) {
	if len(sr.Services) == 0 {
		return nil, false
	}

	findF := func(serAppIds []string) (*ServiceData, bool) {
		for _, appId := range serAppIds {
			service, ok := sr.Services[appId]
			if !ok {
				continue
			}
			if service.ServiceType == proto.ServiceType_ServiceTypeScene && service.MapId != mapId {
				continue
			}
			if !sr.checkAlive(service) {
				serviceLog.Info("remove time service [%v][%v]", service.AppId, service.ServiceType)
				sr.RemoveServiceRecord(appId)
				continue
			}
			return &service, true
		}
		return nil, false
	}

	tarSerAppIds := sr.statusRecord.GetServicesByStatus(ServiceStatusNormal)
	if s, exist = findF(tarSerAppIds); exist {
		return s, exist
	}
	tarSerAppIds = sr.statusRecord.GetServicesByStatus(ServiceStatusFree)
	if s, exist = findF(tarSerAppIds); exist {
		return s, exist
	}
	tarSerAppIds = sr.statusRecord.GetServicesByStatus(ServiceStatusBusy)
	if s, exist = findF(tarSerAppIds); exist {
		return s, exist
	}
	return nil, false
}
