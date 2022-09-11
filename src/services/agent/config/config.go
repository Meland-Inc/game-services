package config

import (
	"fmt"
	"game-message-core/proto"
	"os"
	"strconv"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
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
	StartMs     int64 // 开服时间
}

func (sc *ServiceConfig) Init() error {
	nodeIdStr := os.Getenv("MELAND_SERVICE_AGENT_NODE_ID")
	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil || nodeId == 0 {
		return fmt.Errorf("invalid service id [%v], err: %v", nodeIdStr, err)
	}

	sc.ServiceType = proto.ServiceType_ServiceTypeAgent
	sc.ServerId = nodeId
	sc.StartMs = time_helper.NowUTCMill()
	sc.ServerName = os.Getenv("MELAND_SERVICE_AGENT_DAPR_APPID")
	if sc.ServerName == "" {
		return fmt.Errorf("server app id is empty")
	}

	serviceLog.Info(
		"serviceId [%d],  serviceName [%s],  serviceType [%v]",
		sc.ServerId, sc.ServerName, sc.ServiceType,
	)

	return nil
}
