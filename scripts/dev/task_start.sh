#!/bin/bash
set -o errexit

# ------ export meland config DB settings --------
. scripts/dev/global.sh 

## ---------------meland task service settings ----------------
# taskService         appPort:(5400~5449)  grpc:(5450~5499)   
export APP_ID=game-service-task
export APP_PORT=5400
export DAPR_GRPC_PORT=5450
# export DAPR_HTTP_PORT=5425 
# export APP_API_TOKEN=    
 
echo "---------------------------start DAPR and MAIN service --------------------------------"
dapr run --app-id ${APP_ID} --app-protocol grpc \
--app-port ${APP_PORT} \
--dapr-grpc-port ${DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/task/main.go