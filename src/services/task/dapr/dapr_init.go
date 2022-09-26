package daprService

import (
	"os"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	taskDaprCalls "github.com/Meland-Inc/game-services/src/services/task/dapr/calls"
	taskDaprEvent "github.com/Meland-Inc/game-services/src/services/task/dapr/event"
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
	grpcPort := os.Getenv("MELAND_SERVICE_TASK_DAPR_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = os.Getenv("DAPR_GRPC_PORT")
	}
	serviceLog.Info("dapr grpc port: [%s]", grpcPort)
	return daprInvoke.InitClient(grpcPort)
}

func initDaprService() (err error) {
	appPort := os.Getenv("MELAND_SERVICE_TASK_DAPR_APP_PORT")
	serviceLog.Info("dapr app port: [%s]", appPort)
	if err = daprInvoke.InitServer(appPort); err != nil {
		return err
	}
	if err = taskDaprEvent.InitDaprPubsubEvent(); err != nil {
		return err
	}

	if err = taskDaprCalls.InitDaprCallHandle(); err != nil {
		return err
	}
	return err
}
