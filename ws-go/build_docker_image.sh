#!/bin/bash

DEFAULT_VERSION="latest"

usage() {
  echo '''Usage: ...

  build_image.sh [-flags] [imageVersion]

  -b  build image
  -s  save image to tar.gz
  -p  push image to registry

Example:

  ./build_image.sh  -bp 1.0.0
  '''
}


SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )


BUILD=0
SAVE_IMAGE_TO_TARGZ=0
PUSH_IMAGE=0


args=`getopt bsp $*`
# you should not use `getopt abo: "$@"` since that would parse
# the arguments differently from what the set command below does.
if [[ $? != 0 ]]
then
  usage
  exit 2
fi
set -- ${args}
# You cannot use the set command with a backquoted getopt directly,
# since the exit code from getopt would be shadowed by those of set,
# which is zero by definition.
for i do
  case "$i" in
     -b)
       BUILD=1
       shift;;
     -s)
       SAVE_IMAGE_TO_TARGZ=1
       shift;;
     -p)
       PUSH_IMAGE=1
       shift;;
  esac
done

# Identifying build version
VERSION=$2
if [[ -z "${VERSION}" ]]; then
    VERSION=${DEFAULT_VERSION}
    echo "Version not specified: building with default version [${VERSION}]"
else
    echo "Using specified version [${VERSION}]"
fi


IMAGE_NAME=registry.u-hopper.com/tapoi/tsundoku-blockchain-ws:${VERSION}

cd ${SCRIPT_DIR}

if [[ ${BUILD} == 1 ]]; then
    docker build -t ${IMAGE_NAME} .
    if [[ $? != 0 ]]; then
      echo "Build failed"
      exit 1
    fi

fi

if [[ ${PUSH_IMAGE} == 1 ]]; then

    docker push ${IMAGE_NAME}
    if [[ $? != 0 ]]; then
      echo "Push failed"
      exit 1
    fi
fi

if [[ ${SAVE_IMAGE_TO_TARGZ} == 1 ]]; then
    docker save ${IMAGE_NAME} | gzip > docker_image.tar.gz
fi

if [[ ${BUILD} == 0 ]] && [[ ${PUSH_IMAGE} == 0 ]] && [[ ${SAVE_IMAGE_TO_TARGZ} == 0 ]]; then
  echo "Need to specify at least one parameter (-b, -p, -s)"
  usage
  exit 1
fi