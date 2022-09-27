#!/bin/bash
set -o errexit

## ---------------meland agent service settings ----------------
export MELAND_SERVICE_AGENT_NODE_ID=601
export MELAND_SERVICE_AGENT_DAPR_APPID=meland_service_agent_${MELAND_SERVICE_AGENT_NODE_ID}
export MELAND_SERVICE_AGENT_DAPR_APP_PORT=5600
export MELAND_SERVICE_AGENT_DAPR_GRPC_PORT=5650 
export MELAND_SERVICE_AGENT_SOCKET_HOST=192.168.50.171 
export MELAND_SERVICE_AGENT_SOCKET_PORT=5700 
export MELAND_SERVICE_AGENT_ONLINE_LIMIT=5000


echo "---------------------------start DAPR and agent service --------------------------------"
dapr run --app-id ${MELAND_SERVICE_AGENT_DAPR_APPID} --app-protocol grpc \
--app-port ${MELAND_SERVICE_AGENT_DAPR_APP_PORT} \
--dapr-grpc-port ${MELAND_SERVICE_AGENT_DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/agent/main.go