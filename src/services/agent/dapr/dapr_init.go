package daprService

import (
	"os"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	agentDaprCalls "github.com/Meland-Inc/game-services/src/services/agent/dapr/calls"
	agentDaprEvent "github.com/Meland-Inc/game-services/src/services/agent/dapr/event"
)

func Init() (err error) {
	return initDaprService()
}

func initDaprService() (err error) {
	appPort := os.Getenv("MELAND_SERVICE_AGENT_DAPR_APP_PORT")
	serviceLog.Info("dapr app port: [%s]", appPort)
	if err = daprInvoke.InitServer(appPort); err != nil {
		return err
	}
	if err = agentDaprEvent.InitDaprPubsubEvent(); err != nil {
		return err
	}

	if err = agentDaprCalls.InitDaprCallHandle(); err != nil {
		return err
	}
	return err
}

func Run(errChan chan error) {
	go func() {
		errChan <- daprInvoke.Start()
	}()

	if err := initDaprClient(); err != nil {
		serviceLog.Error("initDaprClient fail err:%v", err)
		panic(err)
	}
}

func initDaprClient() error {
	time.Sleep(time.Millisecond * 300) //300Ms wait dapr link over

	grpcPort := os.Getenv("MELAND_SERVICE_AGENT_DAPR_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = os.Getenv("DAPR_GRPC_PORT")
	}
	serviceLog.Info("dapr grpc port: [%s]", grpcPort)
	return daprInvoke.InitClient(grpcPort)
}
