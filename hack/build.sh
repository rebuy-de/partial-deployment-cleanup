#/bin/bash

cd $( dirname $0 )/..
set -ex

VERSION=$( git describe --always --dirty | tr '-' '.' )

hack/test.sh

mkdir -p target
go build \
	-ldflags "-X main.version=${VERSION}" \
	-o target/partial-deployment-cleanup
