#!/bin/bash

CURDIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)"
BUILDPATH="$(dirname "$CURDIR")/bin"
BUILDTAGS=${BUILDTAGS:-""}

BINNAME=$1
MAINPATH=$2
PLATFORM=$3
SYSPLATFORM=1
OMMITID="$(git rev-parse --short HEAD 2>/dev/null || echo 'na')"
BUILDENV="$(hostname | awk '{print substr($0, 0, 15)}')"
if [[ "$GITLAB_CI" ]]
then
    BUILDENV=gitlabci
fi

help() {
    echo 'Usage: build_binary.sh binary_name binary_main_path'
}

build() {
    local binname=$1
    local mainpath=$2
    local platform=$3



    if [[ "$BUILDTAGS" ]]
    then
        echo "with build tags: ${BUILDTAGS}"
    fi

    local full_name_amd64="${binname}-${platform}-amd64"
    if [ ${platform} == "windows" ]
    	then
    		full_name_amd64="${binname}-${platform}-amd64.exe"
    fi
    local full_path_amd64="${BUILDPATH}/${full_name_amd64}"
    rm -rf "$full_path_amd64"
    echo "building ${full_path_amd64} with commit ${COMMITID} at ${BUILDENV}"
    GOOS="${platform}" GOARCH=amd64 CGO_ENABLED=0 go build \
        -tags "build_bindata ${BUILDTAGS}" \
        -ldflags '-extldflags "-static"  -X bian.CommitId=${COMMITID}  -X bian.BuildEnv=${BUILDENV}'  \
        -o "${full_path_amd64}" \
        "./src/${mainpath}"

}

if [[ -z $BINNAME ]]
then
    help
    exit -1
fi

if [[ -z $MAINPATH ]]
then
    help
    exit -1
fi


case "$OSTYPE" in
    linux*)
        SYSPLATFORM=linux
        ;;
    darwin*)
        SYSPLATFORM=darwin
        ;;
    *)
        echo "unsupported platform $OSTYPE"
        exit -1
        ;;
esac

if [ "$PLATFORM" != "" ]
then
	if [ "$SYSPLATFORM" != "$PLATFORM" ]
	then
		SYSPLATFORM=$PLATFORM
	fi
fi

build "$BINNAME" "$MAINPATH" "$SYSPLATFORM";
