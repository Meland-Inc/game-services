#!/bin/bash
set -o errexit

# ------ export meland config DB settings --------
# . start/global.sh 

## ---------------meland services manager settings  ----------------
export MELAND_SERVICE_MGR_NODE_ID=101
export MELAND_SERVICE_MGR_DAPR_APPID=meland_service_manager
export MELAND_SERVICE_MGR_DAPR_APP_PORT=5100
export MELAND_SERVICE_MGR_DAPR_GRPC_PORT=5150
export MELAND_SERVICE_MGR_HTTP_PORT=5180
 
echo "---------------------------start DAPR and manager service --------------------------------"
dapr run --app-id ${MELAND_SERVICE_MGR_DAPR_APPID} --app-protocol grpc \
--app-port ${MELAND_SERVICE_MGR_DAPR_APP_PORT} \
--dapr-grpc-port ${MELAND_SERVICE_MGR_DAPR_GRPC_PORT} \
--log-level debug -- \
go run ./src/services/manager/main.go 