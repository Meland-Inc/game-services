#!/bin/bash
set -o errexit


## --------------- game services manager settings  ----------------
## appPort:(5100~5129)  grpc:(5130~5169) serviceHttp(5170~5199) 
export APP_ID=game-service-manager
export APP_PORT=5100
export DAPR_GRPC_PORT=5130   
# export DAPR_HTTP_PORT=5125 
# export APP_API_TOKEN= 
export HTTP_PORT=5180
 
echo "---------------------------start DAPR and manager service --------------------------------"
dapr run --app-id ${APP_ID} --app-protocol grpc \
--app-port ${APP_PORT} \
--dapr-grpc-port ${DAPR_GRPC_PORT} \
--log-level debug -- \
go run ./src/services/manager/main.go 