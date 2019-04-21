#!/usr/bin/env bash
#source env.sh or . env.sh

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
MGDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

#export GOPATH=$MGDIR/_gopath
#export PATH=$PATH:$GOPATH/bin