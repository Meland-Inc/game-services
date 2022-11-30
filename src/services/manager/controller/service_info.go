package controller

import (
	base_data "game-message-core/grpc/baseData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

const SERVICE_TIME_OUT_MS int64 = 1000 * 5 // 10 seconds is services timeout

type ServiceData struct {
	AppId           string                    `json:"appId"`
	ServiceType     proto.ServiceType         `json:"serviceType"`
	SceneSerSubType proto.SceneServiceSubType `json:"sceneSerSubType"`
	OwnerId         int64                     `json:"ownerId"`
	Host            string                    `json:"host"`
	Port            int32                     `json:"port"`
	MapId           int32                     `json:"mapId"`
	Online          int32                     `json:"online"`
	MaxOnline       int32                     `json:"maxOnline"`
	CreateAt        int64                     `json:"createdAt"`
	UpdateAt        int64                     `json:"updatedAt"`
}

func (s *ServiceData) ToGrpcService() base_data.ServiceData {
	return base_data.ServiceData{
		AppId:           s.AppId,
		ServiceType:     s.ServiceType,
		SceneSerSubType: s.SceneSerSubType,
		Owner:           s.OwnerId,
		Host:            s.Host,
		Port:            s.Port,
		MapId:           s.MapId,
		Online:          s.Online,
		MaxOnline:       s.MaxOnline,
		CreatedAt:       s.CreateAt,
		UpdatedAt:       s.UpdateAt,
	}
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
	service, exist := sr.Services[appId]
	if !exist {
		return
	}

	delete(sr.Services, appId)
	if service.SceneSerSubType != proto.SceneServiceSubType_Home &&
		service.SceneSerSubType == proto.SceneServiceSubType_Dungeon {
		sr.statusRecord.RemoveServiceStatusRecord(appId)
	}
}

func (sr *ServiceRecord) AddServiceRecord(service ServiceData) bool {
	sr.Services[service.AppId] = service

	if service.SceneSerSubType != proto.SceneServiceSubType_Home &&
		service.SceneSerSubType == proto.SceneServiceSubType_Dungeon {
		sr.statusRecord.AddServiceStatusRecord(service.AppId, service.Online, service.MaxOnline)
	}
	return true
}

func (sr *ServiceRecord) UpdateOrAddServiceRecord(service ServiceData) {
	sr.AddServiceRecord(service)
}

func (sr *ServiceRecord) GetAliveService(
	mapId int32,
	sceneSubType proto.SceneServiceSubType,
	ownerId int64,
) (s *ServiceData, exist bool) {
	if sr.ServiceType != proto.ServiceType_ServiceTypeScene {
		return sr.getPublicAliveService(mapId)
	}
	if sceneSubType == proto.SceneServiceSubType_World {
		return sr.getPublicAliveService(mapId)
	}
	return sr.getPrivateAliveService(mapId, sceneSubType, ownerId)
}

/**
* 公共服务为所有玩家都可以进入的服务,
* 可以动态横向扩展 如大世界地图
* 需要考虑负载均衡 和当前服务压力
**/
func (sr *ServiceRecord) getPublicAliveService(mapId int32) (s *ServiceData, exist bool) {
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

/**
* 玩家专属的服务 在使用完成后会及时回收的服务
* 如(家园 | 副本) 服务
* ownerId 可以是 (玩家id | 队伍id)
**/
func (sr *ServiceRecord) getPrivateAliveService(
	mapId int32,
	sceneSubType proto.SceneServiceSubType,
	ownerId int64,
) (s *ServiceData, exist bool) {
	if len(sr.Services) == 0 {
		return nil, false
	}

	invalidServices := []string{}
	for _, service := range sr.Services {
		if mapId != service.MapId {
			continue
		}
		if service.SceneSerSubType != sceneSubType {
			continue
		}
		if service.OwnerId != ownerId {
			continue
		}
		if !sr.checkAlive(service) {
			serviceLog.Info(
				"remove time service[%v][%v][%v][%v]",
				service.AppId, service.ServiceType, service.SceneSerSubType, service.OwnerId,
			)
			invalidServices = append(invalidServices, service.AppId)
			continue
		}
		s, exist = &service, true
		break
	}

	for _, appId := range invalidServices {
		sr.RemoveServiceRecord(appId)
	}

	return
}
