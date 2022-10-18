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
	serviceLog.Info("received get agentService remote addr:  %s", r.RemoteAddr)
	service, exist := controller.GetInstance().GetAliveServiceByType(proto.ServiceType_ServiceTypeAgent, 0)
	resp := httpData.AgentServiceResp{}
	if !exist {
		resp.ErrorCode = 6000 // TODO: need use global error code
		resp.ErrorMessage = "agent service not found"
	} else {
		resp.Host = service.Host
		resp.Port = service.Port
		resp.Online = service.Online
		resp.MaxOnline = service.MaxOnline
		resp.CreatedAt = service.CreateAt
		resp.UpdateAt = service.UpdateAt
	}
	serviceLog.Info("received get agentService response: %+v ", resp)
	byteArr, err := json.Marshal(resp)
	if err != nil {
		serviceLog.Error("get agent service marshal err: %v", err)
		w.WriteHeader(401)
	}
	w.Write(byteArr)
}
func AllServicesHandler(w http.ResponseWriter, r *http.Request) {
	services := controller.GetInstance().AllServices()
	serviceLog.Info("received get serviceList remote addr: %v, resp: %+v", r.RemoteAddr, services)

	for i := 0; i < len(services)-1; i++ {
		for j := i + 1; j < len(services); j++ {
			if services[i].CreateAt > services[j].CreateAt {
				services[i], services[j] = services[j], services[i]
			}
		}
	}

	byteArr, err := json.Marshal(services)
	if err != nil {
		serviceLog.Error("get all service marshal err: %v", err)
		w.WriteHeader(401)
		return
	}
	w.Write(byteArr)
}
