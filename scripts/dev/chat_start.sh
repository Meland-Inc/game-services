#!/bin/bash
set -o errexit

# ------ export game DB settings --------
. scripts/dev/global.sh 

## ---------------game chat service settings ----------------
## chatService         appPort:(5500~5549)  grpc:(5550~5599)
export APP_ID=game-service-chat
export APP_PORT=5500
export DAPR_GRPC_PORT=5550   
# export DAPR_HTTP_PORT=5525 
# export APP_API_TOKEN= 

echo "---------------------------start DAPR and chat service --------------------------------"
dapr run --app-id ${APP_ID} --app-protocol grpc \
--app-port ${APP_PORT} \
--dapr-grpc-port ${DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/chat/main.go