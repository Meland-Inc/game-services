#!/bin/bash
set -o errexit

# ------ export game DB settings --------
. scripts/dev/global.sh 

## ---------------game account service settings ---------------- 
## appPort:(5200~5249)  grpc:(5250~5299)
export APP_ID=game-service-account
export APP_PORT=5200
export DAPR_GRPC_PORT=5250   
# export DAPR_HTTP_PORT=5225 
# export APP_API_TOKEN= 
 
echo "---------------------------start DAPR and ACCOUNT service --------------------------------"
dapr run --app-id ${APP_ID} --app-protocol grpc \
--app-port ${APP_PORT} \
--dapr-grpc-port ${DAPR_GRPC_PORT} \
--log-level debug -- \
go run src/services/account/main.go