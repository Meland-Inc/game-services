package serviceCnf

import (
	"game-message-core/proto"
)

var instance *ServiceConfig

func GetInstance() *ServiceConfig {
	if instance == nil {
		NewServiceConfig()
	}
	return instance
}

func NewServiceConfig() *ServiceConfig {
	instance = &ServiceConfig{}
	return instance
}

type ServiceConfig struct {
	ServerId    int64
	ServerName  string
	ServiceType proto.ServiceType
	MapId       int32
	Host        string
	Port        int32
	MaxOnline   int32
	StartMs     int64 // 开服时间
}
