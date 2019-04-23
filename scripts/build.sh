#!/bin/bash
set -e
SOURCE=$0;
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
SDIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# Static compile
export CGO_ENABLED=1
# Compile for a determined architecture
export GARCH=amd64

# Compile for a determined host SO
if [[ "$OSTYPE" == "linux-gnu" ]]; then
    export GOOS=linux
elif [[ "$OSTYPE" == "darwin"* ]]; then
    export GOOS=linux
elif [[ "$OSTYPE" == "cygwin" ]]; then
    export GOOS=windows
elif [[ "$OSTYPE" == "msys" ]]; then
    export GOOS=windows
elif [[ "$OSTYPE" == "win32" ]]; then
    export GOOS=windows
elif [[ "$OSTYPE" == "android" ]]; then
    export GOOS=android
elif [[ "$OSTYPE" == "freebsd"* ]]; then
    export GOOS=linux
else
    export GOOS=linux
fi
echo "Compiling for $OSTYPE host"

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