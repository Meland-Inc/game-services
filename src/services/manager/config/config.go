package config

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
	ServerId int64
}

func (sc *ServiceConfig) Init() error {

	return nil
}
