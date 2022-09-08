package config

import (
	"fmt"

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
	fmt.Println("this is demo ---------- service config init -----")
	sc.ServerId = 99999
	sc.ServerName = "demo"
	sc.StartMs = time_helper.NowUTCMill()
	return nil
}
