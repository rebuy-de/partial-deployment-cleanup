#/bin/bash

cd $( dirname $0 )/..
set -ex

hack/build.sh

cd target

rm -rf ./usr/
mkdir -p ./usr/bin/

VERSION=$( ./partial-deployment-cleanup -v )

cp partial-deployment-cleanup ./usr/bin/

fpm \
	-s dir \
	-t rpm \
	-n "rebuy-partial-deployment-cleanup" \
	-v ${VERSION} \
	./usr/bin/partial-deployment-cleanup
