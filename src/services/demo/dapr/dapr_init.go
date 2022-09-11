package daprService

import (
	"os"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	demoDaprCalls "github.com/Meland-Inc/game-services/src/services/demo/dapr/calls"
	demoDaprEvent "github.com/Meland-Inc/game-services/src/services/demo/dapr/event"
)

func Init() (err error) {
	if err = initDaprClient(); err != nil {
		return err
	}

	if err = initDaprService(); err != nil {
		return err
	}

	return nil
}

func initDaprClient() error {
	grpcPort := "5700"
	if grpcPort == "" {
		grpcPort = os.Getenv("DAPR_GRPC_PORT")
	}
	serviceLog.Info("dapr grpc port: [%s]", grpcPort)
	return daprInvoke.InitClient(grpcPort)
}

func initDaprService() (err error) {
	appPort := "5770"
	serviceLog.Info("dapr app port: [%s]", appPort)
	if err = daprInvoke.InitServer(appPort); err != nil {
		return err
	}
	if err = demoDaprEvent.InitDaprPubsubEvent(); err != nil {
		return err
	}

	if err = demoDaprCalls.InitDaprCallHandle(); err != nil {
		return err
	}
	return err
}

func Run() error {
	return daprInvoke.Start()
}

func Stop() {
	daprInvoke.Stop()
}
