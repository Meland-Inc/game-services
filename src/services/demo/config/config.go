package config

import "fmt"

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
}

func (sc *ServiceConfig) Init() error {
	fmt.Println("this is demo ---------- service config init -----")

	return nil
}
