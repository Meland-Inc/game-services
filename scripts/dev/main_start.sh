#!/bin/bash
set -o errexit

# ------ export meland config DB settings --------
. start/global.sh 

## ---------------meland agent service settings ----------------
## mainService         appPort:(5300~5349)  grpc:(5350~5399) 
export MELAND_SERVICE_MAIN_NODE_ID=301
export MELAND_SERVICE_MAIN_DAPR_APPID=meland_service_main
export MELAND_SERVICE_MAIN_DAPR_APP_PORT=5300
export MELAND_SERVICE_MAIN_DAPR_GRPC_PORT=5350
export MELAND_SERVICE_MAIN_DEVELOP_MODEL=true
 
echo "---------------------------start DAPR and MAIN service --------------------------------"
dapr run --app-id ${MELAND_SERVICE_MAIN_DAPR_APPID} --app-protocol grpc \
--app-port ${MELAND_SERVICE_MAIN_DAPR_APP_PORT} \
--dapr-grpc-port ${MELAND_SERVICE_MAIN_DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/main/main.go