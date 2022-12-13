package controller

import (
	base_data "game-message-core/grpc/baseData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

const (
	SER_HEART_TIMEOUT_MS int64 = 1000 * 5
	SER_USER_TIMEOUT_MS  int64 = 1000 * 60 * 3 // user leave service 3 minutes, close user private service
)

type ServiceData struct {
	AppId            string                    `json:"appId"`
	ServiceType      proto.ServiceType         `json:"serviceType"`
	SceneSerSubType  proto.SceneServiceSubType `json:"sceneSerSubType"`
	OwnerId          int64                     `json:"ownerId"`
	Host             string                    `json:"host"`
	Port             int32                     `json:"port"`
	MapId            int32                     `json:"mapId"`
	Online           int32                     `json:"online"`
	MaxOnline        int32                     `json:"maxOnline"`
	CreateAt         int64                     `json:"createdAt"`
	UpdateAt         int64                     `json:"updatedAt"`
	UserLastOnlineAt int64                     `json:"userLastOnlineAt"`
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

func IsUserPrivateSer(ser ServiceData) bool {
	if ser.ServiceType != proto.ServiceType_ServiceTypeScene {
		return false
	}
	if ser.SceneSerSubType == proto.SceneServiceSubType_Home ||
		ser.SceneSerSubType == proto.SceneServiceSubType_Dungeon {
		return true
	}
	return false
}

func CheckAlive(ser ServiceData) bool {
	nowMs := time_helper.NowUTCMill()

	if IsUserPrivateSer(ser) {
		// 非常开服务存活检测
		return nowMs < ser.UserLastOnlineAt+SER_USER_TIMEOUT_MS
	} else {
		return nowMs < ser.UpdateAt+SER_HEART_TIMEOUT_MS
	}
}

func (sr *ServiceRecord) RemoveServiceRecord(appId string) {
	ser, exist := sr.Services[appId]
	if !exist {
		return
	}

	delete(sr.Services, appId)
	if !IsUserPrivateSer(ser) {
		sr.statusRecord.RemoveServiceStatusRecord(appId)
	}
}

func (sr *ServiceRecord) AddServiceRecord(service ServiceData) {
	service.UpdateAt = time_helper.NowUTCMill()
	if service.CreateAt == 0 {
		service.CreateAt = service.UpdateAt
	}
	if service.Online > 0 {
		service.UserLastOnlineAt = service.UpdateAt
	}

	sr.Services[service.AppId] = service
	if !IsUserPrivateSer(service) {
		sr.statusRecord.AddServiceStatusRecord(service.AppId, service.Online, service.MaxOnline)
	}
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
		if mapId != service.MapId ||
			service.SceneSerSubType != sceneSubType ||
			service.OwnerId != ownerId {
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

func (sr *ServiceRecord) checkAndRemoveTimeoutSer(curMs int64) {
	timeOutSerArr := []ServiceData{}
	for _, service := range sr.Services {
		if !CheckAlive(service) {
			timeOutSerArr = append(timeOutSerArr, service)
		}
	}
	for _, ser := range timeOutSerArr {
		serviceLog.Debug(
			"remove time service[%v][%v][%v][%v]",
			ser.AppId, ser.ServiceType, ser.SceneSerSubType, ser.OwnerId,
		)
		sr.RemoveServiceRecord(ser.AppId)
		if IsUserPrivateSer(ser) {
			closeUserPrivateService(ser)
		}
	}
}
