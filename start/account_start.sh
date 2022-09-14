#!/bin/bash
set -o errexit

# ------ export meland config DB settings --------
. start/global.sh 

## ---------------meland agent service settings ----------------
export JWT_SECRET=token key
export MELAND_SERVICE_ACCOUNT_NODE_ID=201
export MELAND_SERVICE_ACCOUNT_DAPR_APPID=meland_service_account
export MELAND_SERVICE_ACCOUNT_DAPR_APP_PORT=5200
export MELAND_SERVICE_ACCOUNT_DAPR_GRPC_PORT=5250  

export MELAND_ACCOUNT_DB_HOST=127.0.0.1
export MELAND_ACCOUNT_DB_USER=root
export MELAND_ACCOUNT_DB_PASS=123456
export MELAND_ACCOUNT_DB_PORT=3306
export MELAND_ACCOUNT_DB_DATABASE=account_dev
 
echo "---------------------------start DAPR and agent service --------------------------------"
dapr run --app-id ${MELAND_SERVICE_ACCOUNT_DAPR_APPID} --app-protocol grpc \
--app-port ${MELAND_SERVICE_ACCOUNT_DAPR_APP_PORT} \
--dapr-grpc-port ${MELAND_SERVICE_ACCOUNT_DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/account/main.go