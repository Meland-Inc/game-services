package httpSer

import (
	"encoding/json"
	"fmt"
	"game-message-core/httpData"
	"game-message-core/proto"
	"net/http"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ping success %v, remote address[%s]", time_helper.NowUTC(), r.RemoteAddr)
}

func AgentServiceHandler(w http.ResponseWriter, r *http.Request) {
	service, exist := controller.GetInstance().GetAliveServiceByType(proto.ServiceType_ServiceTypeGateway)
	resp := httpData.AgentServiceResp{}
	if !exist {
		resp.ErrorCode = 6000 // TODO: need use global error code
		resp.ErrorMessage = "agent service not found"
	} else {
		resp.Host = service.Host
		resp.Port = service.Port
		resp.Online = service.Online
		resp.MaxOnline = service.MaxOnline
		resp.CreatedAt = service.CreatedAt
	}
	byteArr, err := json.Marshal(resp)
	if err != nil {
		serviceLog.Error("get agent service marshal err: %v", err)
	}
	w.Write(byteArr)
}
func AllServicesHandler(w http.ResponseWriter, r *http.Request) {
	services := controller.GetInstance().AllServices()
	byteArr, err := json.Marshal(services)
	if err != nil {
		serviceLog.Error("get agent service marshal err: %v", err)
	}
	w.Write(byteArr)
}
