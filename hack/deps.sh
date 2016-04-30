#!/bin/bash

cd $( dirname $0 )/..
set -ex

hack/consulw.sh version
hack/glidew.sh install
