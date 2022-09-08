package daprService

import (
	"os"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	mgrDaprCalls "github.com/Meland-Inc/game-services/src/services/manager/dapr/calls"
	mgrDaprEvent "github.com/Meland-Inc/game-services/src/services/manager/dapr/event"
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
	grpcPort := os.Getenv("DAPR_GRPC_PORT")
	serviceLog.Info("dapr grpc port: [%s]", grpcPort)
	return daprInvoke.InitClient(grpcPort)
}

func initDaprService() (err error) {
	appPort := os.Getenv("MELAND_SERVICE_MGR_DAPR_APP_PORT")
	serviceLog.Info("dapr app port: [%s]", appPort)
	if err = daprInvoke.InitServer(appPort); err != nil {
		return err
	}
	if err = mgrDaprEvent.InitDaprPubsubEvent(); err != nil {
		return err
	}

	if err = mgrDaprCalls.InitDaprCallHandle(); err != nil {
		return err
	}
	return err
}
