#!/bin/bash
set -o errexit

## ---------------game agent service settings ---------------- 
##  appPort:(5600~5649)  grpc:(5650~5699) serviceSocket(5700~5799)
export APP_ID=game-service-agent-600
export APP_PORT=5600
export DAPR_GRPC_PORT=5650   
# export DAPR_HTTP_PORT=5625 
# export APP_API_TOKEN= 
export SOCKET_HOST=192.168.50.171
export SOCKET_PORT=5700 
export ONLINE_LIMIT=5000


echo "---------------------------start DAPR and agent service --------------------------------" 
dapr run --app-id ${APP_ID} --app-protocol grpc \
--app-port ${APP_PORT} \
--dapr-grpc-port ${DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/agent/main.go