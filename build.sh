#!/bin/bash

# FOLDER_PWD=$(pwd)
SCRIPT=$(realpath "$0")
FOLDER_SCRIPT=$(dirname $SCRIPT)

INTERPRETER=$(ps h -p $$ -o args='' | cut -f1 -d' ' | rev | cut -d'/' -f1 | rev)
if [ $INTERPRETER != "bash" ]; then
    echo "This script should be run with bash, not $interpreter"
    exit 1
fi
if [ -z "${DVS_SRC:-}" ]; then
    echo "DVS_SRC not set"
    exit 1
fi
# Exit More safety, by turning some bugs into errors.
set -o errexit -o pipefail -o nounset

cd $FOLDER_SCRIPT

git checkout chew
git pull

mkdir -p $DVS_SRC/tools/darwin/arm64
CGO_ENABLED=0 GOOS=darwin GOARH=arm64 go build -o $DVS_SRC/tools/darwin/arm64/chew .
mkdir -p $DVS_SRC/tools/linux/amd64
CGO_ENABLED=0 GOOS=linux GOARH=amd64 go build -o $DVS_SRC/tools/linux/amd64/chew .
