#!/bin/bash
set -o errexit

# ------ export meland config DB settings --------
. scripts/dev/global.sh 

## ---------------meland agent service settings ----------------
export MELAND_SERVICE_ACCOUNT_NODE_ID=201
export MELAND_SERVICE_ACCOUNT_DAPR_APPID=meland-service-account
export MELAND_SERVICE_ACCOUNT_DAPR_APP_PORT=5200
export MELAND_SERVICE_ACCOUNT_DAPR_GRPC_PORT=5250  
 
echo "---------------------------start DAPR and ACCOUNT service --------------------------------"
dapr run --app-id ${MELAND_SERVICE_ACCOUNT_DAPR_APPID} --app-protocol grpc \
--app-port ${MELAND_SERVICE_ACCOUNT_DAPR_APP_PORT} \
--dapr-grpc-port ${MELAND_SERVICE_ACCOUNT_DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/account/main.go