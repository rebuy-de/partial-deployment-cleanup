#!/bin/bash

cd $( dirname $0 )/..
set -ex

hack/deps.sh
go test -p 1 $(hack/glidew.sh novendor)
