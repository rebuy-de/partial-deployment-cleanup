#!/bin/bash

cd $( dirname $0 )/..
set -ex

hack/deps.sh

export PATH="$(readlink -f target/consul):${PATH}"
go test -p 1 $(hack/glidew.sh novendor)
