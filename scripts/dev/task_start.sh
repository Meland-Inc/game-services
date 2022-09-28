#!/bin/bash
set -o errexit

# ------ export meland config DB settings --------
. scripts/dev/global.sh 

## ---------------meland task service settings ----------------
# taskService         appPort:(5400~5449)  grpc:(5450~5499)  
export MELAND_SERVICE_TASK_NODE_ID=401
export MELAND_SERVICE_TASK_DAPR_APPID=meland-service-task
export MELAND_SERVICE_TASK_DAPR_APP_PORT=5400
export MELAND_SERVICE_TASK_DAPR_GRPC_PORT=5450 
 
echo "---------------------------start DAPR and MAIN service --------------------------------"
dapr run --app-id ${MELAND_SERVICE_TASK_DAPR_APPID} --app-protocol grpc \
--app-port ${MELAND_SERVICE_TASK_DAPR_APP_PORT} \
--dapr-grpc-port ${MELAND_SERVICE_TASK_DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/task/main.go