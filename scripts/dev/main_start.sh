#!/bin/bash
set -o errexit

# ------ export game DB settings --------
. scripts/dev/global.sh 

## ---------------meland main service settings ----------------
## mainService         appPort:(5300~5349)  grpc:(5350~5399)  
export APP_ID=game-service-main
export APP_PORT=5300
export DAPR_GRPC_PORT=5350   
# export DAPR_HTTP_PORT=5325 
# export APP_API_TOKEN= 
export DEVELOP_MODEL=true
 
echo "---------------------------start DAPR and MAIN service --------------------------------"
dapr run --app-id ${APP_ID} --app-protocol grpc \
--app-port ${APP_PORT} \
--dapr-grpc-port ${DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/main/main.go