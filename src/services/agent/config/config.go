package config

import (
	"fmt"
	"game-message-core/proto"
	"os"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/spf13/cast"
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
	Host        string
	Port        int32
	MaxOnline   int32
	StartMs     int64 // 开服时间
}

func (sc *ServiceConfig) Init() error {
	sc.ServerId = cast.ToInt64(os.Getenv("MELAND_SERVICE_AGENT_NODE_ID"))
	sc.ServiceType = proto.ServiceType_ServiceTypeAgent
	sc.StartMs = time_helper.NowUTCMill()
	sc.ServerName = os.Getenv("MELAND_SERVICE_AGENT_DAPR_APPID")
	sc.Host = os.Getenv("MELAND_SERVICE_AGENT_SOCKET_HOST")
	sc.Port = cast.ToInt32(os.Getenv("MELAND_SERVICE_AGENT_SOCKET_PORT"))
	sc.MaxOnline = cast.ToInt32(os.Getenv("MELAND_SERVICE_AGENT_ONLINE_LIMIT"))
	if sc.MaxOnline == 0 {
		sc.MaxOnline = 5000
	}

	serviceLog.Info(
		"serviceId:[%d], serviceName:[%s], serviceType:[%v], Socket:[%s:%d], maxOnline:[%d]",
		sc.ServerId, sc.ServerName, sc.ServiceType, sc.Host, sc.Port, sc.MaxOnline,
	)

	if sc.ServerId == 0 {
		return fmt.Errorf("invalid serviceId [%v]", sc.ServerId)
	}
	if sc.ServerName == "" {
		return fmt.Errorf("server app id is empty")
	}
	if sc.Port == 0 || sc.Host == "" {
		return fmt.Errorf("invalid socket data, host[%v], port[%v]", sc.Host, sc.Port)
	}
	return nil
}
