set -o errexit
set -o nounset
set -o pipefail

### get project dir
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do
  DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
readonly PROJECT_ROOT="$(dirname $(dirname $DIR))"
RUN_ROOT="$PROJECT_ROOT"
cd $PROJECT_ROOT;

export $(cat "${PROJECT_ROOT}/.env" | xargs)

start=$(date +%s)
platform=linux
BUILDTAGS=${BUILDTAGS:-""}
readonly COMMITID="$(git rev-parse --short HEAD 2>/dev/null || echo 'na')"
readonly BUILDENV="$(hostname | awk '{print substr($0, 0, 15)}')"
readonly DOCKER_NAMESPACE=${DOCKER_NAMESPACE}
readonly IMAGE_TAG=${IMAGE_TAG:-$(date +%Y%m%d%H%M%S)}
readonly SERVICE_TYPE='game-services'

# 输出统一的版本号
omnibus_version=${1:-""}
if ! [ $omnibus_version ]; then
    echo "请输入要部署的环境：👇";
    read -r omnibus_version
fi

export ICW_OMNIBUS_RELEASE="${omnibus_version}"
export ICW_OMNIBUS_VERSION=${IMAGE_TAG}

get_current_version() {
    echo ${IMAGE_TAG}
}

get_docker_tag() {
    appname=${1:-""}
    version=${2:-""}
    echo "${DOCKER_NAMESPACE}/${appname}:${version}"
}

buildGame() {
    local PROJECT_NAME=${SERVICE_TYPE}
    local IMAGE_NAME=$(get_docker_tag $PROJECT_NAME $(get_current_version))
     docker build \
     -f $PROJECT_ROOT/docker/Dockerfile \
     -t $IMAGE_NAME .

    pushImageToRepo $IMAGE_NAME $PROJECT_NAME
}

pushImageToRepo() {
     local IMAGE_NAME=$1
     local PROJECT_NAME=$2
     # 因为image repository 禁止了同一个标签重复写入;
     # 如果push CI 事件和 tag CI 事情使用同一个 commit.sha 去构建 image 会产生一个中断错误;
     # 需要把该类型的错误忽略掉
     if [[ ! $(docker push -q $IMAGE_NAME 2> ${PROJECT_NAME}err.log ) ]]; then
         echo "docker push raw error:$(cat ${PROJECT_NAME}err.log)"
     fi

     if [[  -s ${PROJECT_NAME}err.log  && ! $(cat ${PROJECT_NAME}err.log | grep 'cannot be overwritten because the repository is immutable') ]]; then
        echo "Other errors needs exit"
        exit 1
     fi
}

buildGame

end=$(date +%s)
take=$(( end - start ))
echo "本次构建image.tag = ${IMAGE_TAG}"
echo "✅ ✅ ✅ build done ......... ${take} s"

# 向omnibus-store-record 发送本次build记录
echo "发送本次build记录：" $(curl "https://${OMNIBUS_HOST}/deploy?env=${omnibus_version}&version=${IMAGE_TAG}&name=${SERVICE_TYPE}" -s)