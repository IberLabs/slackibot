#!/bin/bash
SOURCE=$0;
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
MGDIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"