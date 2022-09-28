#!/bin/bash
set -o errexit

# ------ export meland config DB settings --------
. scripts/dev/global.sh 

## ---------------meland chat service settings ----------------
## chatService         appPort:(5500~5549)  grpc:(5550~5599)
export MELAND_SERVICE_CHAT_NODE_ID=501
export MELAND_SERVICE_CHAT_DAPR_APPID=meland-service-chat
export MELAND_SERVICE_CHAT_DAPR_APP_PORT=5500
export MELAND_SERVICE_CHAT_DAPR_GRPC_PORT=5550   
 
echo "---------------------------start DAPR and MAIN service --------------------------------"
dapr run --app-id ${MELAND_SERVICE_CHAT_DAPR_APPID} --app-protocol grpc \
--app-port ${MELAND_SERVICE_CHAT_DAPR_APP_PORT} \
--dapr-grpc-port ${MELAND_SERVICE_CHAT_DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/chat/main.go