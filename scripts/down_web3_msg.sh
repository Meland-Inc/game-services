#!/bin/bash
set -o errexit
# 运行cmcli

### get project dir
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
readonly PROJECT_ROOT="$(dirname $DIR)"
RUN_ROOT="$PROJECT_ROOT"
cd $PROJECT_ROOT;

p=$(curl https://meland-inc.github.io/service-message/message.go)
echo "${p/package main/package message}" > "${PROJECT_ROOT}/src/global/web3Message/message.go"

