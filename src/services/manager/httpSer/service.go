package httpSer

import (
	"fmt"
	"game-message-core/httpData"
	"net/http"
	"os"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func Init() error {
	httpPort := os.Getenv("MELAND_SERVICE_MGR_HTTP_PORT")
	if httpPort == "" {
		err := fmt.Errorf("invalid http port: %v", httpPort)
		serviceLog.Error(err.Error())
		return err
	}

	registerGetHandlers()
	return nil
}

func Run() error {
	addr := fmt.Sprintf(":%s", os.Getenv("MELAND_SERVICE_MGR_HTTP_PORT"))
	serviceLog.Info("start http service address [%s]", addr)
	return http.ListenAndServe(addr, nil)
}

func registerGetHandlers() {
	http.HandleFunc("/ping", PingHandler)
	http.HandleFunc("/serviceList", AllServicesHandler)
	http.HandleFunc("/"+httpData.API_GET_AGENT_SERVICE, AgentServiceHandler)
}
