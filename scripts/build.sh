#!/usr/bin/env bash

set -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

export CGO_ENABLED=1

source ${SDIR}/env.sh
MGDIR_GOPATH=${MGDIR}

pushd ${MGDIR_GOPATH}/cmd/slackibot

#BUILD_TIME="$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
#TAG="current"
#REVISION="current"
#if hash git 2>/dev/null && [ -e ${MGDIR_GOPATH}/.git ]; then
#  TAG="$(git describe --tags)"
#  REVISION="$(git rev-parse HEAD)"
#fi


go build
mv slackibot* ${MGDIR_GOPATH}/bin/
popd
#rm -rfv ${MGDIR_GOPATH}/bin