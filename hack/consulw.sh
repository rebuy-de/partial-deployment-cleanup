#!/bin/bash

cd $( dirname $0 )/..
set -e

if [ ! -f target/consul/consul ]
then
	set -x
	mkdir -p target/consul

	VERSION=0.6.4
	FILE=consul_${VERSION}_$(go env GOHOSTOS)_$(go env GOHOSTARCH).zip
	BASE=https://releases.hashicorp.com/consul
	URL=${BASE}/${VERSION}/${FILE}

	wget -c \
		-O target/consul/${FILE} \
		${URL}
	unzip \
		target/consul/${FILE} \
		-d target/consul

	set +x
fi

target/consul/consul "$@"
