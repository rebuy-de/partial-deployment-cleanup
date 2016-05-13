#/bin/bash

source $( dirname $0)/build.sh

cd target

rm -rf ./usr/
mkdir -p ./usr/bin/

cp ${BINARY_NAME} ./usr/bin/partial-deployment-cleanup

fpm \
	-s dir \
	-t rpm \
	-n "rebuy-partial-deployment-cleanup" \
	-v ${VERSION} \
	./usr/bin/partial-deployment-cleanup
