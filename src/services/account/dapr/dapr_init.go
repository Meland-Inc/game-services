package daprService

import (
	"os"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	accountDaprCalls "github.com/Meland-Inc/game-services/src/services/account/dapr/calls"
	accountDaprEvent "github.com/Meland-Inc/game-services/src/services/account/dapr/event"
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
	grpcPort := os.Getenv("MELAND_SERVICE_ACCOUNT_DAPR_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = os.Getenv("DAPR_GRPC_PORT")
	}
	serviceLog.Info("dapr grpc port: [%s]", grpcPort)
	return daprInvoke.InitClient(grpcPort)
}

func initDaprService() (err error) {
	appPort := os.Getenv("MELAND_SERVICE_ACCOUNT_DAPR_APP_PORT")
	serviceLog.Info("dapr app port: [%s]", appPort)
	if err = daprInvoke.InitServer(appPort); err != nil {
		return err
	}
	if err = accountDaprEvent.InitDaprPubsubEvent(); err != nil {
		return err
	}

	if err = accountDaprCalls.InitDaprCallHandle(); err != nil {
		return err
	}
	return err
}
