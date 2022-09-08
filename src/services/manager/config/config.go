package config

import (
	"fmt"
	"os"

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
	ServerId   int64
	ServerName string
	StartMs    int64 // 开服时间
}

func (sc *ServiceConfig) Init() error {
	sc.ServerId = 1001
	sc.StartMs = time_helper.NowUTCMill()
	sc.ServerName = os.Getenv("MELAND_SERVICE_MGR_DAPR_APPID")
	if sc.ServerName == "" {
		return fmt.Errorf("server app id is empty")
	}
	return nil
}
