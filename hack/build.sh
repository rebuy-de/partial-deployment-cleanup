#/bin/bash

cd $( dirname $0 )/..
set -ex

hack/test.sh
go build
