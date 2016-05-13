#!/bin/bash

cd $( dirname $0 )/..
set -ex

export VERSION=$( git describe --always --dirty | tr '-' '.' )
export BINARY_NAME=partial-deployment-cleanup-${VERSION}

export PATH="$(readlink -f target/consul):${PATH}"

hack/consulw.sh version
hack/glidew.sh install
